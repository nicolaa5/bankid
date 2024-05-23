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

<img width=800 src="./authflow.jpg"/>

### Test Setup
> [!IMPORTANT]  
> The personnummer you use for the Test BankID has to be valid. See the [following list for Personnummers](https://github.com/emilybache/personnummer/blob/master/valid_100.txt) that are valid

1. Set up Mobile BankID on your phone (Android/iOS) or the BankID Security Application on your computer with a [test configuration](https://www.bankid.com/en/utvecklare/test/skaffa-testbankid/testbankid-konfiguration)
2. Create a [Test BankID](https://www.bankid.com/en/utvecklare/test/skaffa-testbankid/test-bankid-get) at https://demo.bankid.comthat are accepted by BankID
3. Run the CLI program wiht `bankid auth --test` in order to test authentication with your Test BankID

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