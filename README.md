<div align="center">
    <img width=300 src="./bankid-go.png"/>
</div>

# 🇸🇪 BankID
## Relying Party (RP) service
![ Unit Tests](https://github.com/nicolaa5/bankid/actions/workflows/unit.tests.yml/badge.svg)  

### Who is this repository for? 
You can use this repo if you're using BankID in your organization for one of the following purposes: 
- Authenticating users to use your services
- Signing documents, transactions or payments related to your organization

### Setup
For setup instructions from BankID visit [their integration guide](https://www.bankid.com/en/utvecklare/guider/teknisk-integrationsguide/rp-introduktion)

### Examples
```go
// Provide certificate and URL 
b, err := bankid.New(bankid.Parameters{
    URL: bankid.BankIDTestUrl,
    Certificate: bankid.Certificate{
        Passphrase:     bankid.BankIDTestPassphrase,
        SSLCertificate: bankid.SSLTestCertificate,
        CACertificate:  bankid.CATestCertificate,
    },
})

// Send authenticate request to BankID
authResponse, err := b.Auth(bankid.AuthRequest{
    EndUserIP: ip,
    Requirement: bankid.Requirement{
        PersonalNumber: personNummer,
    },
})

// Poll for the status of the order
collectResponse, err = b.Collect(bankid.CollectRequest{
    OrderRef: authResponse.OrderRef,
})
```