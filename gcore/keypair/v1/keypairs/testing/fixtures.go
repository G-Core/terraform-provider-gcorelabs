package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/keypair/v1/keypairs"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaYBMuzTylHSRUi7ZF4KW07sNxd1rHnHzruruoX0ZLSaRMslI0Wp9RV+rvJyfqf98mCFWCfvQxkiaH0pjEoSZkTiLUz/dbLWjeR+K5TAKCcdzfKGA8aeCwdrAFR7rMKeQ2GrVebZCgfEIKCXiCHv0hObujDKCaTFJxg8Q8sXCwmg835bfSB+Wz9QW2zAX6olEttsXEOjG2aaz006WYq4Eo3NEeBizUNpwC3+HifXXcmJhNQtb+3Uf6bb5RVaY2nSRkk2MXmf8eXfakAO0j/ek3o9pVqjGAuTIJdPLpM4cSAV9jP+WNHZ2nu0u0VvjBO2/BLx950sas5DS78OzLpQ9H",
      "sshkey_id": "keypair",
      "sshkey_name": "keypair",
      "fingerprint": "b1:f1:f4:61:a5:a8:eb:6c:7f:4f:e3:a9:d5:bd:4d:46"
    }
  ]
}
`

const GetResponse = `
{
  "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaYBMuzTylHSRUi7ZF4KW07sNxd1rHnHzruruoX0ZLSaRMslI0Wp9RV+rvJyfqf98mCFWCfvQxkiaH0pjEoSZkTiLUz/dbLWjeR+K5TAKCcdzfKGA8aeCwdrAFR7rMKeQ2GrVebZCgfEIKCXiCHv0hObujDKCaTFJxg8Q8sXCwmg835bfSB+Wz9QW2zAX6olEttsXEOjG2aaz006WYq4Eo3NEeBizUNpwC3+HifXXcmJhNQtb+3Uf6bb5RVaY2nSRkk2MXmf8eXfakAO0j/ek3o9pVqjGAuTIJdPLpM4cSAV9jP+WNHZ2nu0u0VvjBO2/BLx950sas5DS78OzLpQ9H",
  "sshkey_id": "keypair",
  "sshkey_name": "keypair",
  "fingerprint": "b1:f1:f4:61:a5:a8:eb:6c:7f:4f:e3:a9:d5:bd:4d:46"
}
`

const CreateRequest = `
{
  "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaYBMuzTylHSRUi7ZF4KW07sNxd1rHnHzruruoX0ZLSaRMslI0Wp9RV+rvJyfqf98mCFWCfvQxkiaH0pjEoSZkTiLUz/dbLWjeR+K5TAKCcdzfKGA8aeCwdrAFR7rMKeQ2GrVebZCgfEIKCXiCHv0hObujDKCaTFJxg8Q8sXCwmg835bfSB+Wz9QW2zAX6olEttsXEOjG2aaz006WYq4Eo3NEeBizUNpwC3+HifXXcmJhNQtb+3Uf6bb5RVaY2nSRkk2MXmf8eXfakAO0j/ek3o9pVqjGAuTIJdPLpM4cSAV9jP+WNHZ2nu0u0VvjBO2/BLx950sas5DS78OzLpQ9H",
  "sshkey_name": "keypair"
}
`

const CreateResponse = `
{
  "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaYBMuzTylHSRUi7ZF4KW07sNxd1rHnHzruruoX0ZLSaRMslI0Wp9RV+rvJyfqf98mCFWCfvQxkiaH0pjEoSZkTiLUz/dbLWjeR+K5TAKCcdzfKGA8aeCwdrAFR7rMKeQ2GrVebZCgfEIKCXiCHv0hObujDKCaTFJxg8Q8sXCwmg835bfSB+Wz9QW2zAX6olEttsXEOjG2aaz006WYq4Eo3NEeBizUNpwC3+HifXXcmJhNQtb+3Uf6bb5RVaY2nSRkk2MXmf8eXfakAO0j/ek3o9pVqjGAuTIJdPLpM4cSAV9jP+WNHZ2nu0u0VvjBO2/BLx950sas5DS78OzLpQ9H",
  "sshkey_id": "keypair",
  "sshkey_name": "keypair",
  "fingerprint": "b1:f1:f4:61:a5:a8:eb:6c:7f:4f:e3:a9:d5:bd:4d:46"
}
`

var (
	KeyPair1 = keypairs.KeyPair{
		Name:        "keypair",
		ID:          "keypair",
		Fingerprint: "b1:f1:f4:61:a5:a8:eb:6c:7f:4f:e3:a9:d5:bd:4d:46",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDaYBMuzTylHSRUi7ZF4KW07sNxd1rHnHzruruoX0ZLSaRMslI0Wp9RV+rvJyfqf98mCFWCfvQxkiaH0pjEoSZkTiLUz/dbLWjeR+K5TAKCcdzfKGA8aeCwdrAFR7rMKeQ2GrVebZCgfEIKCXiCHv0hObujDKCaTFJxg8Q8sXCwmg835bfSB+Wz9QW2zAX6olEttsXEOjG2aaz006WYq4Eo3NEeBizUNpwC3+HifXXcmJhNQtb+3Uf6bb5RVaY2nSRkk2MXmf8eXfakAO0j/ek3o9pVqjGAuTIJdPLpM4cSAV9jP+WNHZ2nu0u0VvjBO2/BLx950sas5DS78OzLpQ9H",
		PrivateKey:  nil,
	}

	ExpectedKeyPairsSlice = []keypairs.KeyPair{KeyPair1}
)
