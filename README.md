# taubyte/domain-validation

This repos is used to validate ownership of DNS domains

# Generate Keys
To use this package you need to generate a symetric pair for keys:

```
openssl ecparam -name prime256v1 -genkey -noout -out private.key
```
```
openssl ec -in tprivate.key -pubout -out public.pem
```

NOTE: the private key need to be kept safe.


# Maintainers
 - Samy Fodil @samyfodil
 - Tafseer Khan @tafseer-khan
 - Aron Jalbuena @arontaubyte
