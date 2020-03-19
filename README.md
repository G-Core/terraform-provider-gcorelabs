Installation
------------------------------------
    
    export GOPRIVATE="bitbucket.gcore.lu"
    VERSION=1.14.1
    curl -Ssqo /tmp/go${VERSION}.linux-amd64.tar.gz https://dl.google.com/go/go${VERSION}.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -xvf /tmp/go${VERSION}.linux-amd64.tar.gz -C /usr/local

    git config --global url."git@bitbucket.gcore.lu:".insteadOf "https://bitbucket.gcore.lu/scm/"
    
    cat <EOF >> ~/.ssh/config    
    Host bitbucket.gcore.lu
      AddKeysToAgent yes
      User git
      Port 7999
      IdentityFile /path/to/key
      GSSAPIAuthentication no
    EOF
