## Certificates
 You have to apply for a SSL certificate for your company by signing a contract<br> with [one of the banks listed here](https://www.bankid.com/en/foretag/kontakt-foeretag).

The following certificates are meant to be used with BankID's [`test`](https://www.bankid.com/en/utvecklare/guider/verification-of-digital-id-card/test-environment) environment
| Certificate Name | Description                                                                                           |
|------------------|--------------------------------------------------------------------------------------------------|
| [ca_test.crt](https://www.bankid.com/en/utvecklare/guider/verification-of-digital-id-card/test-environment)  | **Issuer of server certificate:**<br> CN = Test BankID SSL Root CA v1 Test<br> OU = Infrastructure CA<br> O = Finansiell ID-Teknik BID AB    | 
| [ca_prod.crt](https://www.bankid.com/en/utvecklare/guider/verification-of-digital-id-card/production-environment)  | **Issuer of server certificate:**<br> The server certificate is issued by the following CA.<br> CN = BankID SSL Root CA v1<br> OU = Infrastructure CA<br> O = Finansiell ID-Teknik BID AB |
| [FPTestcert5_20240610.p12](https://www.bankid.com/en/utvecklare/test)   | **Certificate for test**<br> TLS certificate for test<br>                        |


-------

🇸🇪 Swedish explanation of the test certificates: 
**1. FPTestcert5_20240610.p12:**

* Den här filen lagrar certifikatet och den privata nyckeln i PKCS#12-format.
* Den krypteras med AES-256-CBC-algoritmen som har högre säkerhet än äldre metoder.
* När du skapar ditt certifikat för produktion med BankID Keygen kommer det skapas i det här formatet. 

**2. FPTestcert5_20240610.pem:**

* Den här filen innehåller certifikatet och den krypterade privata nyckeln i PEM-format. 
* Certifikatet ligger i början av filen, följt av den privata nyckeln.

