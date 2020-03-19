package testing

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticatedHeaders(t *testing.T) {
	p := &gcorecloud.ProviderClient{
		AccessTokenID: "1234",
	}
	assert.Len(t, p.AccessToken(), len("1234"))
	expected := map[string]string{"Authorization": "Bearer 1234"}
	actual := p.AuthenticatedHeaders()
	assert.Equal(t, expected, actual)
}

func TestUserAgent(t *testing.T) {
	p := &gcorecloud.ProviderClient{}

	p.UserAgent.Prepend("custom-user-agent/2.4.0")
	expected := "custom-user-agent/2.4.0 gcorecloud/0.0.1"
	actual := p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)

	p.UserAgent.Prepend("another-custom-user-agent/0.3.0", "a-third-ua/5.9.0")
	expected = "another-custom-user-agent/0.3.0 a-third-ua/5.9.0 custom-user-agent/2.4.0 gcorecloud/0.0.1"
	actual = p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)

	p.UserAgent = gcorecloud.UserAgent{}
	expected = "gcorecloud/0.0.1"
	actual = p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)
}

func TestConcurrentReauth(t *testing.T) {
	var info = struct {
		numreauths  int
		failedAuths int
		mut         *sync.RWMutex
	}{
		0,
		0,
		new(sync.RWMutex),
	}

	numConc := 20

	atc := client.NewAuthResultTest(client.AccessToken, client.RefreshToken)

	preReAuthAccessToken := client.AccessToken
	postReAuthToken := "12345678"
	postReAuthTokenHeader := fmt.Sprintf("Bearer %s", postReAuthToken)

	p := new(gcorecloud.ProviderClient)
	p.UseTokenLock()
	err := p.SetTokensAndAuthResult(atc)
	if err != nil {
		log.Error(err)
	}
	p.ReauthFunc = func() error {
		p.SetThrowaway(true)
		time.Sleep(1 * time.Second)
		p.AuthenticatedHeaders()
		info.mut.Lock()
		info.numreauths++
		info.mut.Unlock()
		p.AccessTokenID = postReAuthToken
		p.SetThrowaway(false)
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != postReAuthTokenHeader {
			w.WriteHeader(http.StatusUnauthorized)
			info.mut.Lock()
			info.failedAuths++
			info.mut.Unlock()
			return
		}
		info.mut.RLock()
		hasReauthed := info.numreauths != 0
		info.mut.RUnlock()

		if hasReauthed {
			th.CheckEquals(t, p.AccessToken(), postReAuthToken)
		}

		w.Header().Add("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{}`)
		if err != nil {
			log.Error(err)
		}
	})

	wg := new(sync.WaitGroup)
	reqopts := new(gcorecloud.RequestOpts)
	reqopts.MoreHeaders = map[string]string{
		"X-Auth-AccessToken": preReAuthAccessToken,
	}

	for i := 0; i < numConc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := p.Request("GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts)
			th.CheckNoErr(t, err)
			if resp == nil {
				t.Errorf("got a nil response")
				return
			}
			if resp.Body == nil {
				t.Errorf("response body was nil")
				return
			}
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					log.Error(err)
				}
			}()
			actual, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("error reading response body: %s", err)
				return
			}
			th.CheckByteArrayEquals(t, []byte(`{}`), actual)
		}()
	}

	wg.Wait()

	th.AssertEquals(t, 1, info.numreauths)
}

func TestReauthEndLoop(t *testing.T) {
	var info = struct {
		reauthAttempts   int
		maxReauthReached bool
		mut              *sync.RWMutex
	}{
		0,
		false,
		new(sync.RWMutex),
	}

	numconc := 20
	mut := new(sync.RWMutex)
	atc := client.NewAuthResultTest(client.AccessToken, client.RefreshToken)

	p := new(gcorecloud.ProviderClient)
	p.UseTokenLock()
	err := p.SetTokensAndAuthResult(atc)
	if err != nil {
		log.Error(err)
	}
	p.ReauthFunc = func() error {
		info.mut.Lock()
		defer info.mut.Unlock()

		if info.reauthAttempts > 5 {
			info.maxReauthReached = true
			return fmt.Errorf("max reauthentication attempts reached")
		}
		p.SetThrowaway(true)
		p.AuthenticatedHeaders()
		p.SetThrowaway(false)
		info.reauthAttempts++

		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		// route always return 401
		w.WriteHeader(http.StatusUnauthorized)
	})

	reqopts := new(gcorecloud.RequestOpts)

	// counters for the upcoming errors
	errAfter := 0
	errUnable := 0

	wg := new(sync.WaitGroup)
	for i := 0; i < numconc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Request("GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts) // nolint

			mut.Lock()
			defer mut.Unlock()

			// ErrErrorAfter... will happen after a successful reauthentication,
			// but the service still responds with a 401.
			if _, ok := err.(*gcorecloud.ErrErrorAfterReauthentication); ok {
				errAfter++
			}

			// ErrErrorUnable... will happen when the custom reauth func reports
			// an error.
			if _, ok := err.(*gcorecloud.ErrUnableToReauthenticate); ok {
				errUnable++
			}
		}()
	}

	wg.Wait()
	th.AssertEquals(t, info.reauthAttempts, 6)
	th.AssertEquals(t, info.maxReauthReached, true)
	th.AssertEquals(t, errAfter > 1, true)
	th.AssertEquals(t, errUnable < 20, true)
}

func TestRequestThatCameDuringReauthWaitsUntilItIsCompleted(t *testing.T) {
	var info = struct {
		numreauths  int
		failedAuths int
		reauthCh    chan struct{}
		mut         *sync.RWMutex
	}{
		0,
		0,
		make(chan struct{}),
		new(sync.RWMutex),
	}

	numconc := 20

	preReAuthToken := client.AccessToken
	postReAuthToken := "12345678"
	atc := client.NewAuthResultTest(client.AccessToken, client.RefreshToken)
	postReAuthTokenHeader := fmt.Sprintf("Bearer %s", postReAuthToken)

	p := new(gcorecloud.ProviderClient)
	p.UseTokenLock()
	err := p.SetTokensAndAuthResult(atc)
	if err != nil {
		log.Error(err)
	}
	p.ReauthFunc = func() error {
		info.mut.RLock()
		if info.numreauths == 0 {
			info.mut.RUnlock()
			close(info.reauthCh)
			time.Sleep(1 * time.Second)
		} else {
			info.mut.RUnlock()
		}
		p.SetThrowaway(true)
		p.AuthenticatedHeaders()
		info.mut.Lock()
		info.numreauths++
		info.mut.Unlock()
		p.AccessTokenID = postReAuthToken
		p.SetThrowaway(false)
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != postReAuthTokenHeader {
			info.mut.Lock()
			info.failedAuths++
			info.mut.Unlock()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		info.mut.RLock()
		hasReauthed := info.numreauths != 0
		info.mut.RUnlock()

		if hasReauthed {
			th.CheckEquals(t, p.AccessToken(), postReAuthToken)
		}

		w.Header().Add("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{}`)
		if err != nil {
			log.Error(err)
		}
	})

	wg := new(sync.WaitGroup)
	reqopts := new(gcorecloud.RequestOpts)
	reqopts.MoreHeaders = map[string]string{
		"X-Auth-AccessToken": preReAuthToken,
	}

	for i := 0; i < numconc; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i != 0 {
				<-info.reauthCh
			}
			resp, err := p.Request("GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts)
			th.CheckNoErr(t, err)
			if resp == nil {
				t.Errorf("got a nil response")
				return
			}
			if resp.Body == nil {
				t.Errorf("response body was nil")
				return
			}
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					log.Error(err)
				}
			}()
			actual, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("error reading response body: %s", err)
				return
			}
			th.CheckByteArrayEquals(t, []byte(`{}`), actual)
		}(i)
	}

	wg.Wait()

	th.AssertEquals(t, 1, info.numreauths)
	th.AssertEquals(t, 1, info.failedAuths)
}

func TestRequestReauthsAtMostOnce(t *testing.T) {
	// There was an issue where GCore cloud would go into an infinite
	// reauthentication loop with buggy services that send 401 even for fresh
	// tokens. This test simulates such a service and checks that a call to
	// ProviderClient.Request() will not try to reauthenticate more than once.

	reauthCounter := 0
	var reauthCounterMutex sync.Mutex
	atc := client.NewAuthResultTest(client.AccessToken, client.RefreshToken)

	p := new(gcorecloud.ProviderClient)
	p.UseTokenLock()
	err := p.SetTokensAndAuthResult(atc)
	if err != nil {
		log.Error(err)
	}
	p.ReauthFunc = func() error {
		reauthCounterMutex.Lock()
		reauthCounter++
		reauthCounterMutex.Unlock()
		//The actual token value does not matter, the endpoint does not check it.
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	requestCounter := 0
	var requestCounterMutex sync.Mutex

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		requestCounterMutex.Lock()
		requestCounter++
		//avoid infinite loop
		if requestCounter == 10 {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		requestCounterMutex.Unlock()

		//always reply 401, even immediately after reauthenticate
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	// The expected error message indicates that we reauthenticated once (that's
	// the part before the colon), but when encountering another 401 response, we
	// did not attempt reauthentication again and just passed that 401 response to
	// the caller as ErrDefault401.
	_, err = p.Request("GET", th.Endpoint()+"/route", &gcorecloud.RequestOpts{}) // nolint
	expectedErrorMessage := "Successfully re-authenticated, but got error executing request: Authentication failed"
	if err != nil {
		th.AssertEquals(t, expectedErrorMessage, err.Error())
	}
}

func TestRequestWithContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "OK")
		if err != nil {
			log.Error(err)
		}
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	p := &gcorecloud.ProviderClient{Context: ctx}

	res, err := p.Request("GET", ts.URL, &gcorecloud.RequestOpts{})
	th.AssertNoErr(t, err)
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}
	err = res.Body.Close()
	if err != nil {
		log.Error(err)
	}
	th.AssertNoErr(t, err)

	cancel()
	_, err = p.Request("GET", ts.URL, &gcorecloud.RequestOpts{}) // nolint
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	if !strings.Contains(err.Error(), ctx.Err().Error()) {
		t.Fatalf("expecting error to contain: %q, got %q", ctx.Err().Error(), err.Error())
	}
}
