<div align="center">
    <img width=200 src="./bankid-go.png"/>
</div>

# 🇸🇪 BankID
![ Unit Tests](https://github.com/nicolaa5/bankid/actions/workflows/unit.tests.yml/badge.svg)  

## Relying Party (RP) service

To install the CLI app with which to interact with the package
```bash
go install github.com/nicolaa5/bankid/cmd/bankid
```

### Who is this repository for? 
You can use this repo if you're using BankID in your organization for one of the following purposes: 
- Authenticating users to use your services
- Signing documents, transactions or payments related to your organization

<img width=800 src="./authflow.jpg"/>

### Test Setup
> [!IMPORTANT]  
> The personnummer you use for the Test BankID has to be valid. See the [following list for Personnummers](https://github.com/emilybache/personnummer/blob/master/valid_100.txt) that are valid

1. Set up Mobile BankID on your phone (Android/iOS) or the BankID Security Application on your computer with a [test configuration](https://www.bankid.com/en/utvecklare/test/skaffa-testbankid/testbankid-konfiguration)
2. Create a [Test BankID](https://www.bankid.com/en/utvecklare/test/skaffa-testbankid/test-bankid-get) at https://demo.bankid.com that are accepted by BankID
3. Run the CLI program with `bankid auth --test` in order to test authentication with your Test BankID

### Examples
```go
// Provide certificate and URL
b, err := bankid.New(bankid.Config{
    URL: bankid.BankIDURL,
    Certificate: bankid.Certificate{
        Passphrase:     passphrase,
        SSLCertificate: sslCert,
    },
})

// Send authenticate request to BankID
authResponse, err := b.Auth(ctx, bankid.AuthRequest{
    EndUserIP: ip,
    Requirement: &bankid.Requirement{
        PersonalNumber: personNummer,
    },
})

// Poll for the status of the order
collectResponse, err := b.Collect(ctx, bankid.CollectRequest{
    OrderRef: authResponse.OrderRef,
})

fmt.Println(collectResponse)
// Success case
// {
//         "orderRef": "5cc86d87-ded0-43c3-8ce8-7693710a0092",
//         "status": "complete",
//         "completionData": {
//                 "user": {
//                         "personalNumber": "199510221287",
//                         "name": "John Doe",
//                         "givenName": "John",
//                         "surname": "Doe"
//                 },
//                 ...
//         }
//  }
```