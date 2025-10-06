## 
          ![](/static/shared/icon/chrome_reader_mode_gm_grey_24dp.svg) README [¶](#section-readme)


### Package validator


Package validator implements value validations for structs and individual fields based on tags.


It has the following **unique** features: 


              - Cross Field and Cross Struct validations by using validation tags or custom validators.

              - Slice, Array and Map diving, which allows any or all levels of a multidimensional field to be validated.

              - Ability to dive into both map keys and values for validation

              - Handles type interface by determining it's underlying type prior to validation.

              - Handles custom field types such as sql driver Valuer see [Valuer](https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29)


              - Alias validation tags, which allows for mapping of several validations to a single tag for easier defining of validations on structs

              - Extraction of custom defined Field Name e.g. can specify to extract the JSON name while validating and have it available in the resulting FieldError

              - Customizable i18n aware error messages.

              - Default validator for the [gin](https://github.com/gin-gonic/gin) web framework; upgrading from v8 to v9 in gin see [here](https://github.com/go-playground/validator/tree/master/_examples/gin-upgrading-overriding)


#### A Call for Maintainers


Please read the discussiong started [here](https://github.com/go-playground/validator/discussions/1330) if you are interested in contributing/helping maintain this package. 


#### Installation


Use go get.


```

                go get github.com/go-playground/validator/v10


```


Then import the validator package into your own code.


```

                import "github.com/go-playground/validator/v10"


```


#### Error Return Value


Validation functions return type error


They return type error to avoid the issue discussed in the following, where err is always != nil:


              - 
                [http://stackoverflow.com/a/29138676/3158232](http://stackoverflow.com/a/29138676/3158232)


              - 
                [https://github.com/go-playground/validator/issues/134](https://github.com/go-playground/validator/issues/134)


Validator returns only InvalidValidationError for bad validation input, nil or ValidationErrors as type error; so, in your code all you need to do is check if the error returned is not nil, and if it's not check if error is InvalidValidationError ( if necessary, most of the time it isn't ) type cast it to type ValidationErrors like so:


```

                err := validate.Struct(mystruct)
validationErrors := err.(validator.ValidationErrors)


```


#### Usage and documentation


Please see [https://pkg.go.dev/github.com/go-playground/validator/v10](https://pkg.go.dev/github.com/go-playground/validator/v10) for detailed usage docs. 

            Examples:


              - 
                [Simple](https://github.com/go-playground/validator/raw/master/_examples/simple/main.go)


              - 
                [Custom Field Types](https://github.com/go-playground/validator/raw/master/_examples/custom/main.go)


              - 
                [Struct Level](https://github.com/go-playground/validator/raw/master/_examples/struct-level/main.go)


              - 
                [Translations & Custom Errors](https://github.com/go-playground/validator/raw/master/_examples/translations/main.go)


              - 
                [Gin upgrade and/or override validator](https://github.com/go-playground/validator/tree/v9/_examples/gin-upgrading-overriding)


              - 
                [wash - an example application putting it all together](https://github.com/bluesuncorp/wash)


#### Baked-in Validations


##### Special Notes:


              - If new to using validator it is highly recommended to initialize it using the `WithRequiredStructEnabled` option which is opt-in to new behaviour that will become the default behaviour in v11+. See documentation for more details. 


```

                validate := validator.New(validator.WithRequiredStructEnabled())


```


##### Fields:


| Tag | Description |
| --- | --- |
| eqcsfield | Field Equals Another Field (relative) |
| eqfield | Field Equals Another Field |
| fieldcontains | Check the indicated characters are present in the Field |
| fieldexcludes | Check the indicated characters are not present in the field |
| gtcsfield | Field Greater Than Another Relative Field |
| gtecsfield | Field Greater Than or Equal To Another Relative Field |
| gtefield | Field Greater Than or Equal To Another Field |
| gtfield | Field Greater Than Another Field |
| ltcsfield | Less Than Another Relative Field |
| ltecsfield | Less Than or Equal To Another Relative Field |
| ltefield | Less Than or Equal To Another Field |
| ltfield | Less Than Another Field |
| necsfield | Field Does Not Equal Another Field (relative) |
| nefield | Field Does Not Equal Another Field |


##### Network:


| Tag | Description |
| --- | --- |
| cidr | Classless Inter-Domain Routing CIDR |
| cidrv4 | Classless Inter-Domain Routing CIDRv4 |
| cidrv6 | Classless Inter-Domain Routing CIDRv6 |
| datauri | Data URL |
| fqdn | Full Qualified Domain Name (FQDN) |
| hostname | Hostname RFC 952 |
| hostname_port | HostPort |
| hostname_rfc1123 | Hostname RFC 1123 |
| ip | Internet Protocol Address IP |
| ip4_addr | Internet Protocol Address IPv4 |
| ip6_addr | Internet Protocol Address IPv6 |
| ip_addr | Internet Protocol Address IP |
| ipv4 | Internet Protocol Address IPv4 |
| ipv6 | Internet Protocol Address IPv6 |
| mac | Media Access Control Address MAC |
| tcp4_addr | Transmission Control Protocol Address TCPv4 |
| tcp6_addr | Transmission Control Protocol Address TCPv6 |
| tcp_addr | Transmission Control Protocol Address TCP |
| udp4_addr | User Datagram Protocol Address UDPv4 |
| udp6_addr | User Datagram Protocol Address UDPv6 |
| udp_addr | User Datagram Protocol Address UDP |
| unix_addr | Unix domain socket end point Address |
| uri | URI String |
| url | URL String |
| http_url | HTTP URL String |
| url_encoded | URL Encoded |
| urn_rfc2141 | Urn RFC 2141 String |


##### Strings:


| Tag | Description |
| --- | --- |
| alpha | Alpha Only |
| alphanum | Alphanumeric |
| alphanumunicode | Alphanumeric Unicode |
| alphaunicode | Alpha Unicode |
| ascii | ASCII |
| boolean | Boolean |
| contains | Contains |
| containsany | Contains Any |
| containsrune | Contains Rune |
| endsnotwith | Ends Not With |
| endswith | Ends With |
| excludes | Excludes |
| excludesall | Excludes All |
| excludesrune | Excludes Rune |
| lowercase | Lowercase |
| multibyte | Multi-Byte Characters |
| number | Number |
| numeric | Numeric |
| printascii | Printable ASCII |
| startsnotwith | Starts Not With |
| startswith | Starts With |
| uppercase | Uppercase |


##### Format:


| Tag | Description |
| --- | --- |
| base64 | Base64 String |
| base64url | Base64URL String |
| base64rawurl | Base64RawURL String |
| bic | Business Identifier Code (ISO 9362) |
| bcp47_language_tag | Language tag (BCP 47) |
| btc_addr | Bitcoin Address |
| btc_addr_bech32 | Bitcoin Bech32 Address (segwit) |
| credit_card | Credit Card Number |
| mongodb | MongoDB ObjectID |
| mongodb_connection_string | MongoDB Connection String |
| cron | Cron |
| spicedb | SpiceDb ObjectID/Permission/Type |
| datetime | Datetime |
| e164 | e164 formatted phone number |
| ein | U.S. Employeer Identification Number |
| email | E-mail String |
| eth_addr | Ethereum Address |
| hexadecimal | Hexadecimal String |
| hexcolor | Hexcolor String |
| hsl | HSL String |
| hsla | HSLA String |
| html | HTML Tags |
| html_encoded | HTML Encoded |
| isbn | International Standard Book Number |
| isbn10 | International Standard Book Number 10 |
| isbn13 | International Standard Book Number 13 |
| issn | International Standard Serial Number |
| iso3166_1_alpha2 | Two-letter country code (ISO 3166-1 alpha-2) |
| iso3166_1_alpha3 | Three-letter country code (ISO 3166-1 alpha-3) |
| iso3166_1_alpha_numeric | Numeric country code (ISO 3166-1 numeric) |
| iso3166_2 | Country subdivision code (ISO 3166-2) |
| iso4217 | Currency code (ISO 4217) |
| json | JSON |
| jwt | JSON Web Token (JWT) |
| latitude | Latitude |
| longitude | Longitude |
| luhn_checksum | Luhn Algorithm Checksum (for strings and (u)int) |
| postcode_iso3166_alpha2 | Postcode |
| postcode_iso3166_alpha2_field | Postcode |
| rgb | RGB String |
| rgba | RGBA String |
| ssn | Social Security Number SSN |
| timezone | Timezone |
| uuid | Universally Unique Identifier UUID |
| uuid3 | Universally Unique Identifier UUID v3 |
| uuid3_rfc4122 | Universally Unique Identifier UUID v3 RFC4122 |
| uuid4 | Universally Unique Identifier UUID v4 |
| uuid4_rfc4122 | Universally Unique Identifier UUID v4 RFC4122 |
| uuid5 | Universally Unique Identifier UUID v5 |
| uuid5_rfc4122 | Universally Unique Identifier UUID v5 RFC4122 |
| uuid_rfc4122 | Universally Unique Identifier UUID RFC4122 |
| md4 | MD4 hash |
| md5 | MD5 hash |
| sha256 | SHA256 hash |
| sha384 | SHA384 hash |
| sha512 | SHA512 hash |
| ripemd128 | RIPEMD-128 hash |
| ripemd128 | RIPEMD-160 hash |
| tiger128 | TIGER128 hash |
| tiger160 | TIGER160 hash |
| tiger192 | TIGER192 hash |
| semver | Semantic Versioning 2.0.0 |
| ulid | Universally Unique Lexicographically Sortable Identifier ULID |
| cve | Common Vulnerabilities and Exposures Identifier (CVE id) |


##### Comparisons:


| Tag | Description |
| --- | --- |
| eq | Equals |
| eq_ignore_case | Equals ignoring case |
| gt | Greater than |
| gte | Greater than or equal |
| lt | Less Than |
| lte | Less Than or Equal |
| ne | Not Equal |
| ne_ignore_case | Not Equal ignoring case |


##### Other:


                  ``


| Tag | Description |
| --- | --- |
| dir | Existing Directory |
| dirpath | Directory Path |
| file | Existing File |
| filepath | File Path |
| image | Image |
| isdefault | Is Default |
| len | Length |
| max | Maximum |
| min | Minimum |
| oneof | One Of |
| required | Required |
| required_if | Required If |
| required_unless | Required Unless |
| required_with | Required With |
| required_with_all | Required With All |
| required_without | Required Without |
| required_without_all | Required Without All |
| excluded_if | Excluded If |
| excluded_unless | Excluded Unless |
| excluded_with | Excluded With |
| excluded_with_all | Excluded With All |
| excluded_without | Excluded Without |
| excluded_without_all | Excluded Without All |
| unique | Unique |
| validateFn | Verify if the method Validate() error does not return an error (or any specified method) |


###### Aliases:


| Tag | Description |
| --- | --- |
| iscolor | hexcolor|rgb|rgba|hsl|hsla |
| country_code | iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric |


#### Benchmarks

            Run on MacBook Pro Max M3

```

                go version go1.23.3 darwin/arm64
goos: darwin
goarch: arm64
cpu: Apple M3 Max
pkg: github.com/go-playground/validator/v10
BenchmarkFieldSuccess-16                                                42461943                27.88 ns/op            0 B/op          0 allocs/op
BenchmarkFieldSuccessParallel-16                                        486632887                2.289 ns/op           0 B/op          0 allocs/op
BenchmarkFieldFailure-16                                                 9566167               121.3 ns/op           200 B/op          4 allocs/op
BenchmarkFieldFailureParallel-16                                        17551471                83.68 ns/op          200 B/op          4 allocs/op
BenchmarkFieldArrayDiveSuccess-16                                        7602306               155.6 ns/op            97 B/op          5 allocs/op
BenchmarkFieldArrayDiveSuccessParallel-16                               20664610                59.80 ns/op           97 B/op          5 allocs/op
BenchmarkFieldArrayDiveFailure-16                                        4659756               252.9 ns/op           301 B/op         10 allocs/op
BenchmarkFieldArrayDiveFailureParallel-16                                8010116               152.9 ns/op           301 B/op         10 allocs/op
BenchmarkFieldMapDiveSuccess-16                                          2834575               421.2 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveSuccessParallel-16                                  7179700               171.8 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveFailure-16                                          3081728               384.4 ns/op           376 B/op         13 allocs/op
BenchmarkFieldMapDiveFailureParallel-16                                  6058137               204.0 ns/op           377 B/op         13 allocs/op
BenchmarkFieldMapDiveWithKeysSuccess-16                                  2544975               464.8 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveWithKeysSuccessParallel-16                          6661954               181.4 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveWithKeysFailure-16                                  2435484               490.7 ns/op           553 B/op         16 allocs/op
BenchmarkFieldMapDiveWithKeysFailureParallel-16                          4249617               282.0 ns/op           554 B/op         16 allocs/op
BenchmarkFieldCustomTypeSuccess-16                                      14943525                77.35 ns/op           32 B/op          2 allocs/op
BenchmarkFieldCustomTypeSuccessParallel-16                              64051954                20.61 ns/op           32 B/op          2 allocs/op
BenchmarkFieldCustomTypeFailure-16                                      10721384               107.1 ns/op           184 B/op          3 allocs/op
BenchmarkFieldCustomTypeFailureParallel-16                              18714495                69.77 ns/op          184 B/op          3 allocs/op
BenchmarkFieldOrTagSuccess-16                                            4063124               294.3 ns/op            16 B/op          1 allocs/op
BenchmarkFieldOrTagSuccessParallel-16                                   31903756                41.22 ns/op           18 B/op          1 allocs/op
BenchmarkFieldOrTagFailure-16                                            7748558               146.8 ns/op           216 B/op          5 allocs/op
BenchmarkFieldOrTagFailureParallel-16                                   13139854                92.05 ns/op          216 B/op          5 allocs/op
BenchmarkStructLevelValidationSuccess-16                                16808389                70.25 ns/op           16 B/op          1 allocs/op
BenchmarkStructLevelValidationSuccessParallel-16                        90686955                14.47 ns/op           16 B/op          1 allocs/op
BenchmarkStructLevelValidationFailure-16                                 5818791               200.2 ns/op           264 B/op          7 allocs/op
BenchmarkStructLevelValidationFailureParallel-16                        11115874               107.5 ns/op           264 B/op          7 allocs/op
BenchmarkStructSimpleCustomTypeSuccess-16                                7764956               151.9 ns/op            32 B/op          2 allocs/op
BenchmarkStructSimpleCustomTypeSuccessParallel-16                       52316265                30.37 ns/op           32 B/op          2 allocs/op
BenchmarkStructSimpleCustomTypeFailure-16                                4195429               277.2 ns/op           416 B/op          9 allocs/op
BenchmarkStructSimpleCustomTypeFailureParallel-16                        7305661               164.6 ns/op           432 B/op         10 allocs/op
BenchmarkStructFilteredSuccess-16                                        6312625               186.1 ns/op           216 B/op          5 allocs/op
BenchmarkStructFilteredSuccessParallel-16                               13684459                93.42 ns/op          216 B/op          5 allocs/op
BenchmarkStructFilteredFailure-16                                        6751482               171.2 ns/op           216 B/op          5 allocs/op
BenchmarkStructFilteredFailureParallel-16                               14146070                86.93 ns/op          216 B/op          5 allocs/op
BenchmarkStructPartialSuccess-16                                         6544448               177.3 ns/op           224 B/op          4 allocs/op
BenchmarkStructPartialSuccessParallel-16                                13951946                88.73 ns/op          224 B/op          4 allocs/op
BenchmarkStructPartialFailure-16                                         4075833               287.5 ns/op           440 B/op          9 allocs/op
BenchmarkStructPartialFailureParallel-16                                 7490805               161.3 ns/op           440 B/op          9 allocs/op
BenchmarkStructExceptSuccess-16                                          4107187               281.4 ns/op           424 B/op          8 allocs/op
BenchmarkStructExceptSuccessParallel-16                                 15979173                80.86 ns/op          208 B/op          3 allocs/op
BenchmarkStructExceptFailure-16                                          4434372               264.3 ns/op           424 B/op          8 allocs/op
BenchmarkStructExceptFailureParallel-16                                  8081367               154.1 ns/op           424 B/op          8 allocs/op
BenchmarkStructSimpleCrossFieldSuccess-16                                6459542               183.4 ns/op            56 B/op          3 allocs/op
BenchmarkStructSimpleCrossFieldSuccessParallel-16                       41013781                37.95 ns/op           56 B/op          3 allocs/op
BenchmarkStructSimpleCrossFieldFailure-16                                4034998               292.1 ns/op           272 B/op          8 allocs/op
BenchmarkStructSimpleCrossFieldFailureParallel-16                       11348446               115.3 ns/op           272 B/op          8 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldSuccess-16                     4448528               267.7 ns/op            64 B/op          4 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldSuccessParallel-16            26813619                48.33 ns/op           64 B/op          4 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldFailure-16                     3090646               384.5 ns/op           288 B/op          9 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldFailureParallel-16             9870906               129.5 ns/op           288 B/op          9 allocs/op
BenchmarkStructSimpleSuccess-16                                         10675562               109.5 ns/op             0 B/op          0 allocs/op
BenchmarkStructSimpleSuccessParallel-16                                 131159784                8.932 ns/op           0 B/op          0 allocs/op
BenchmarkStructSimpleFailure-16                                          4094979               286.6 ns/op           416 B/op          9 allocs/op
BenchmarkStructSimpleFailureParallel-16                                  7606663               157.9 ns/op           416 B/op          9 allocs/op
BenchmarkStructComplexSuccess-16                                         2073470               576.0 ns/op           224 B/op          5 allocs/op
BenchmarkStructComplexSuccessParallel-16                                 7821831               161.3 ns/op           224 B/op          5 allocs/op
BenchmarkStructComplexFailure-16                                          576358              2001 ns/op            3042 B/op         48 allocs/op
BenchmarkStructComplexFailureParallel-16                                 1000000              1171 ns/op            3041 B/op         48 allocs/op
BenchmarkOneof-16                                                       22503973                52.82 ns/op            0 B/op          0 allocs/op
BenchmarkOneofParallel-16                                                8538474               140.4 ns/op             0 B/op          0 allocs/op


```


#### Complementary Software


Here is a list of software that complements using this library either pre or post validation.


              - 
                [form](https://github.com/go-playground/form) - Decodes url.Values into Go value(s) and Encodes Go value(s) into url.Values. Dual Array and Full map support.


              - 
                [mold](https://github.com/go-playground/mold) - A general library to help modify or set data within data structures and other objects


#### How to Contribute


Make a pull request...


#### Maintenance and support for SDK major versions


See prior discussion [here](https://github.com/go-playground/validator/discussions/1342) for more details. 


This package is aligned with the [Go release policy](https://go.dev/doc/devel/release) in that support is guaranteed for the two most recent major versions. 


This does not mean the package will not work with older versions of Go, only that we reserve the right to increase the MSGV(Minimum Supported Go Version) when the need arises to address Security issues/patches, OS issues & support or newly introduced functionality that would greatly benefit the maintenance and/or usage of this package.


If and when the MSGV is increased it will be done so in a minimum of a `Minor` release bump. 


#### License


Distributed under MIT License, please see license file within the code for more details.


#### Maintainers


This project has grown large enough that more than one person is required to properly support the community. If you are interested in becoming a maintainer please reach out to me [https://github.com/deankarn](https://github.com/deankarn)


        Expand ▾
        Collapse ▴


## 
          ![](/static/shared/icon/code_gm_grey_24dp.svg) Documentation [¶](#section-documentation)


### Overview [¶](#pkg-overview)


                  - 
                    [Singleton](#hdr-Singleton)


                  - 
                    [Validation Functions Return Type error](#hdr-Validation_Functions_Return_Type_error)


                  - 
                    [Custom Validation Functions](#hdr-Custom_Validation_Functions)


                  - 
                    [Cross-Field Validation](#hdr-Cross_Field_Validation)


                  - 
                    [Multiple Validators](#hdr-Multiple_Validators)


                  - 
                    [Using Validator Tags](#hdr-Using_Validator_Tags)


                  - 
                    [Baked In Validators and Tags](#hdr-Baked_In_Validators_and_Tags)


                  - 
                    [Skip Field](#hdr-Skip_Field)


                  - 
                    [Or Operator](#hdr-Or_Operator)


                  - 
                    [StructOnly](#hdr-StructOnly)


                  - 
                    [NoStructLevel](#hdr-NoStructLevel)


                  - 
                    [Omit Empty](#hdr-Omit_Empty)


                  - 
                    [Omit Nil](#hdr-Omit_Nil)


                  - 
                    [Dive](#hdr-Dive)


                  - 
                    [Required](#hdr-Required)


                  - 
                    [Required If](#hdr-Required_If)


                  - 
                    [Required Unless](#hdr-Required_Unless)


                  - 
                    [Required With](#hdr-Required_With)


                  - 
                    [Required With All](#hdr-Required_With_All)


                  - 
                    [Required Without](#hdr-Required_Without)


                  - 
                    [Required Without All](#hdr-Required_Without_All)


                  - 
                    [Excluded If](#hdr-Excluded_If)


                  - 
                    [Excluded Unless](#hdr-Excluded_Unless)


                  - 
                    [Is Default](#hdr-Is_Default)


                  - 
                    [Length](#hdr-Length)


                  - 
                    [Maximum](#hdr-Maximum)


                  - 
                    [Minimum](#hdr-Minimum)


                  - 
                    [Equals](#hdr-Equals)


                  - 
                    [Not Equal](#hdr-Not_Equal)


                  - 
                    [One Of](#hdr-One_Of)


                  - 
                    [One Of Case Insensitive](#hdr-One_Of_Case_Insensitive)


                  - 
                    [Greater Than](#hdr-Greater_Than)


                  - 
                    [Greater Than or Equal](#hdr-Greater_Than_or_Equal)


                  - 
                    [Less Than](#hdr-Less_Than)


                  - 
                    [Less Than or Equal](#hdr-Less_Than_or_Equal)


                  - 
                    [Field Equals Another Field](#hdr-Field_Equals_Another_Field)


                  - 
                    [Field Does Not Equal Another Field](#hdr-Field_Does_Not_Equal_Another_Field)


                  - 
                    [Field Greater Than Another Field](#hdr-Field_Greater_Than_Another_Field)


                  - 
                    [Field Greater Than Another Relative Field](#hdr-Field_Greater_Than_Another_Relative_Field)


                  - 
                    [Field Greater Than or Equal To Another Field](#hdr-Field_Greater_Than_or_Equal_To_Another_Field)


                  - 
                    [Field Greater Than or Equal To Another Relative Field](#hdr-Field_Greater_Than_or_Equal_To_Another_Relative_Field)


                  - 
                    [Less Than Another Field](#hdr-Less_Than_Another_Field)


                  - 
                    [Less Than Another Relative Field](#hdr-Less_Than_Another_Relative_Field)


                  - 
                    [Less Than or Equal To Another Field](#hdr-Less_Than_or_Equal_To_Another_Field)


                  - 
                    [Less Than or Equal To Another Relative Field](#hdr-Less_Than_or_Equal_To_Another_Relative_Field)


                  - 
                    [Field Contains Another Field](#hdr-Field_Contains_Another_Field)


                  - 
                    [Field Excludes Another Field](#hdr-Field_Excludes_Another_Field)


                  - 
                    [Unique](#hdr-Unique)


                  - 
                    [ValidateFn](#hdr-ValidateFn)


                  - 
                    [Alpha Only](#hdr-Alpha_Only)


                  - 
                    [Alphanumeric](#hdr-Alphanumeric)


                  - 
                    [Alpha Unicode](#hdr-Alpha_Unicode)


                  - 
                    [Alphanumeric Unicode](#hdr-Alphanumeric_Unicode)


                  - 
                    [Boolean](#hdr-Boolean)


                  - 
                    [Number](#hdr-Number)


                  - 
                    [Numeric](#hdr-Numeric)


                  - 
                    [Hexadecimal String](#hdr-Hexadecimal_String)


                  - 
                    [Hexcolor String](#hdr-Hexcolor_String)


                  - 
                    [Lowercase String](#hdr-Lowercase_String)


                  - 
                    [Uppercase String](#hdr-Uppercase_String)


                  - 
                    [RGB String](#hdr-RGB_String)


                  - 
                    [RGBA String](#hdr-RGBA_String)


                  - 
                    [HSL String](#hdr-HSL_String)


                  - 
                    [HSLA String](#hdr-HSLA_String)


                  - 
                    [E.164 Phone Number String](#hdr-E_164_Phone_Number_String)


                  - 
                    [E-mail String](#hdr-E_mail_String)


                  - 
                    [JSON String](#hdr-JSON_String)


                  - 
                    [JWT String](#hdr-JWT_String)


                  - 
                    [File](#hdr-File)


                  - 
                    [Image path](#hdr-Image_path)


                  - 
                    [File Path](#hdr-File_Path)


                  - 
                    [URL String](#hdr-URL_String)


                  - 
                    [URI String](#hdr-URI_String)


                  - 
                    [Urn](#hdr-Urn_RFC_2141_String)
                    [RFC 2141](https://rfc-editor.org/rfc/rfc2141.html) String


                  - 
                    [Base32 String](#hdr-Base32_String)


                  - 
                    [Base64 String](#hdr-Base64_String)


                  - 
                    [Base64URL String](#hdr-Base64URL_String)


                  - 
                    [Base64RawURL String](#hdr-Base64RawURL_String)


                  - 
                    [Bitcoin Address](#hdr-Bitcoin_Address)


                  - 
                    [Ethereum Address](#hdr-Ethereum_Address)


                  - 
                    [Contains](#hdr-Contains)


                  - 
                    [Contains Any](#hdr-Contains_Any)


                  - 
                    [Contains Rune](#hdr-Contains_Rune)


                  - 
                    [Excludes](#hdr-Excludes)


                  - 
                    [Excludes All](#hdr-Excludes_All)


                  - 
                    [Excludes Rune](#hdr-Excludes_Rune)


                  - 
                    [Starts With](#hdr-Starts_With)


                  - 
                    [Ends With](#hdr-Ends_With)


                  - 
                    [Does Not Start With](#hdr-Does_Not_Start_With)


                  - 
                    [Does Not End With](#hdr-Does_Not_End_With)


                  - 
                    [International Standard Book Number](#hdr-International_Standard_Book_Number)


                  - 
                    [International Standard Book Number 10](#hdr-International_Standard_Book_Number_10)


                  - 
                    [International Standard Book Number 13](#hdr-International_Standard_Book_Number_13)


                  - 
                    [Universally Unique Identifier UUID](#hdr-Universally_Unique_Identifier_UUID)


                  - 
                    [Universally Unique Identifier UUID v3](#hdr-Universally_Unique_Identifier_UUID_v3)


                  - 
                    [Universally Unique Identifier UUID v4](#hdr-Universally_Unique_Identifier_UUID_v4)


                  - 
                    [Universally Unique Identifier UUID v5](#hdr-Universally_Unique_Identifier_UUID_v5)


                  - 
                    [Universally Unique Lexicographically Sortable Identifier ULID](#hdr-Universally_Unique_Lexicographically_Sortable_Identifier_ULID)


                  - 
                    [ASCII](#hdr-ASCII)


                  - 
                    [Printable ASCII](#hdr-Printable_ASCII)


                  - 
                    [Multi-Byte Characters](#hdr-Multi_Byte_Characters)


                  - 
                    [Data URL](#hdr-Data_URL)


                  - 
                    [Latitude](#hdr-Latitude)


                  - 
                    [Longitude](#hdr-Longitude)


                  - 
                    [Employeer Identification Number EIN](#hdr-Employeer_Identification_Number_EIN)


                  - 
                    [Social Security Number SSN](#hdr-Social_Security_Number_SSN)


                  - 
                    [Internet Protocol Address IP](#hdr-Internet_Protocol_Address_IP)


                  - 
                    [Internet Protocol Address IPv4](#hdr-Internet_Protocol_Address_IPv4)


                  - 
                    [Internet Protocol Address IPv6](#hdr-Internet_Protocol_Address_IPv6)


                  - 
                    [Classless Inter-Domain Routing CIDR](#hdr-Classless_Inter_Domain_Routing_CIDR)


                  - 
                    [Classless Inter-Domain Routing CIDRv4](#hdr-Classless_Inter_Domain_Routing_CIDRv4)


                  - 
                    [Classless Inter-Domain Routing CIDRv6](#hdr-Classless_Inter_Domain_Routing_CIDRv6)


                  - 
                    [Transmission Control Protocol Address TCP](#hdr-Transmission_Control_Protocol_Address_TCP)


                  - 
                    [Transmission Control Protocol Address TCPv4](#hdr-Transmission_Control_Protocol_Address_TCPv4)


                  - 
                    [Transmission Control Protocol Address TCPv6](#hdr-Transmission_Control_Protocol_Address_TCPv6)


                  - 
                    [User Datagram Protocol Address UDP](#hdr-User_Datagram_Protocol_Address_UDP)


                  - 
                    [User Datagram Protocol Address UDPv4](#hdr-User_Datagram_Protocol_Address_UDPv4)


                  - 
                    [User Datagram Protocol Address UDPv6](#hdr-User_Datagram_Protocol_Address_UDPv6)


                  - 
                    [Internet Protocol Address IP](#hdr-Internet_Protocol_Address_IP-1)


                  - 
                    [Internet Protocol Address IPv4](#hdr-Internet_Protocol_Address_IPv4-1)


                  - 
                    [Internet Protocol Address IPv6](#hdr-Internet_Protocol_Address_IPv6-1)


                  - 
                    [Unix domain socket end point Address](#hdr-Unix_domain_socket_end_point_Address)


                  - 
                    [Media Access Control Address MAC](#hdr-Media_Access_Control_Address_MAC)


                  - 
                    [Hostname](#hdr-Hostname_RFC_952)
                    [RFC 952](https://rfc-editor.org/rfc/rfc952.html)


                  - 
                    [Hostname](#hdr-Hostname_RFC_1123)
                    [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html)


                  - 
                    [HTML Tags](#hdr-HTML_Tags)


                  - 
                    [HTML Encoded](#hdr-HTML_Encoded)


                  - 
                    [URL Encoded](#hdr-URL_Encoded)


                  - 
                    [Directory](#hdr-Directory)


                  - 
                    [Directory Path](#hdr-Directory_Path)


                  - 
                    [HostPort](#hdr-HostPort)


                  - 
                    [Datetime](#hdr-Datetime)


                  - 
                    [Iso3166-1 alpha-2](#hdr-Iso3166_1_alpha_2)


                  - 
                    [Iso3166-1 alpha-3](#hdr-Iso3166_1_alpha_3)


                  - 
                    [Iso3166-1 alpha-numeric](#hdr-Iso3166_1_alpha_numeric)


                  - 
                    [BCP 47 Language Tag](#hdr-BCP_47_Language_Tag)


                  - 
                    [](#hdr-RFC_1035_label)
                    [RFC 1035](https://rfc-editor.org/rfc/rfc1035.html) label


                  - 
                    [TimeZone](#hdr-TimeZone)


                  - 
                    [Semantic Version](#hdr-Semantic_Version)


                  - 
                    [CVE Identifier](#hdr-CVE_Identifier)


                  - 
                    [Credit Card](#hdr-Credit_Card)


                  - 
                    [Luhn Checksum](#hdr-Luhn_Checksum)


                  - 
                    [MongoDB](#hdr-MongoDB)


                  - 
                    [Cron](#hdr-Cron)


                  - 
                    [SpiceDb ObjectID/Permission/Object Type](#hdr-SpiceDb_ObjectID_Permission_Object_Type)


                  - 
                    [Alias Validators and Tags](#hdr-Alias_Validators_and_Tags)


                  - 
                    [Non standard validators](#hdr-Non_standard_validators)


                  - 
                    [Panics](#hdr-Panics)


Package validator implements value validations for structs and individual fields based on tags. 


It can also handle Cross-Field and Cross-Struct validation for nested structs and has the ability to dive into arrays and maps of any type. 


see more examples [https://github.com/go-playground/validator/tree/master/_examples](https://github.com/go-playground/validator/tree/master/_examples)


#### Singleton [¶](#hdr-Singleton)


Validator is designed to be thread-safe and used as a singleton instance. It caches information about your struct and validations, in essence only parsing your validation tags once per struct type. Using multiple instances neglects the benefit of caching. The not thread-safe functions are explicitly marked as such in the documentation. 


#### Validation Functions Return Type error [¶](#hdr-Validation_Functions_Return_Type_error)


Doing things this way is actually the way the standard library does, see the file.Open method here: 


```
https://golang.org/pkg/os/#Open.
```


The authors return type "error" to avoid the issue discussed in the following, where err is always != nil: 


```
http://stackoverflow.com/a/29138676/3158232
https://github.com/go-playground/validator/issues/134
```


Validator only InvalidValidationError for bad validation input, nil or ValidationErrors as type error; so, in your code all you need to do is check if the error returned is not nil, and if it's not check if error is InvalidValidationError ( if necessary, most of the time it isn't ) type cast it to type ValidationErrors like so err.(validator.ValidationErrors). 


#### Custom Validation Functions [¶](#hdr-Custom_Validation_Functions)


Custom Validation functions can be added. Example: 


```
// Structure
func customFunc(fl validator.FieldLevel) bool {

  if fl.Field().String() == "invalid" {
    return false
  }

  return true
}

validate.RegisterValidation("custom tag name", customFunc)
// NOTES: using the same tag name as an existing function
//        will overwrite the existing one
```


#### Cross-Field Validation [¶](#hdr-Cross_Field_Validation)


Cross-Field Validation can be done via the following tags: 


                - eqfield

                - nefield

                - gtfield

                - gtefield

                - ltfield

                - ltefield

                - eqcsfield

                - necsfield

                - gtcsfield

                - gtecsfield

                - ltcsfield

                - ltecsfield


If, however, some custom cross-field validation is required, it can be done using a custom validation. 


Why not just have cross-fields validation tags (i.e. only eqcsfield and not eqfield)? 


The reason is efficiency. If you want to check a field within the same struct "eqfield" only has to find the field on the same struct (1 level). But, if we used "eqcsfield" it could be multiple levels down. Example: 


```
type Inner struct {
  StartDate time.Time
}

type Outer struct {
  InnerStructField *Inner
  CreatedAt time.Time      `validate:"ltecsfield=InnerStructField.StartDate"`
}

now := time.Now()

inner := &Inner{
  StartDate: now,
}

outer := &Outer{
  InnerStructField: inner,
  CreatedAt: now,
}

errs := validate.Struct(outer)

// NOTE: when calling validate.Struct(val) topStruct will be the top level struct passed
//       into the function
//       when calling validate.VarWithValue(val, field, tag) val will be
//       whatever you pass, struct, field...
//       when calling validate.Field(field, tag) val will be nil
```


#### Multiple Validators [¶](#hdr-Multiple_Validators)


Multiple validators on a field will process in the order defined. Example: 


```
type Test struct {
  Field `validate:"max=10,min=1"`
}

// max will be checked then min
```


Bad Validator definitions are not handled by the library. Example: 


```
type Test struct {
  Field `validate:"min=10,max=0"`
}

// this definition of min max will never succeed
```


#### Using Validator Tags [¶](#hdr-Using_Validator_Tags)


Baked In Cross-Field validation only compares fields on the same struct. If Cross-Field + Cross-Struct validation is needed you should implement your own custom validator. 


Comma (",") is the default separator of validation tags. If you wish to have a comma included within the parameter (i.e. excludesall=,) you will need to use the UTF-8 hex representation 0x2C, which is replaced in the code as a comma, so the above will become excludesall=0x2C. 


```
type Test struct {
  Field `validate:"excludesall=,"`    // BAD! Do not include a comma.
  Field `validate:"excludesall=0x2C"` // GOOD! Use the UTF-8 hex representation.
}
```


Pipe ("|") is the 'or' validation tags deparator. If you wish to have a pipe included within the parameter i.e. excludesall=| you will need to use the UTF-8 hex representation 0x7C, which is replaced in the code as a pipe, so the above will become excludesall=0x7C 


```
type Test struct {
  Field `validate:"excludesall=|"`    // BAD! Do not include a pipe!
  Field `validate:"excludesall=0x7C"` // GOOD! Use the UTF-8 hex representation.
}
```


#### Baked In Validators and Tags [¶](#hdr-Baked_In_Validators_and_Tags)


Here is a list of the current built in validators: 


#### Skip Field [¶](#hdr-Skip_Field)


Tells the validation to skip this struct field; this is particularly handy in ignoring embedded structs from being validated. (Usage: -) 


```
Usage: -
```


#### Or Operator [¶](#hdr-Or_Operator)


This is the 'or' operator allowing multiple validators to be used and accepted. (Usage: rgb|rgba) <-- this would allow either rgb or rgba colors to be accepted. This can also be combined with 'and' for example ( Usage: omitempty,rgb|rgba) 


```
Usage: |
```


#### StructOnly [¶](#hdr-StructOnly)


When a field that is a nested struct is encountered, and contains this flag any validation on the nested struct will be run, but none of the nested struct fields will be validated. This is useful if inside of your program you know the struct will be valid, but need to verify it has been assigned. NOTE: only "required" and "omitempty" can be used on a struct itself. 


```
Usage: structonly
```


#### NoStructLevel [¶](#hdr-NoStructLevel)


Same as structonly tag except that any struct level validations will not run. 


```
Usage: nostructlevel
```


#### Omit Empty [¶](#hdr-Omit_Empty)


Allows conditional validation, for example, if a field is not set with a value (Determined by the "required" validator) then other validation such as min or max won't run, but if a value is set validation will run. 


```
Usage: omitempty
```


#### Omit Nil [¶](#hdr-Omit_Nil)


Allows to skip the validation if the value is nil (same as omitempty, but only for the nil-values). 


```
Usage: omitnil
```


#### Dive [¶](#hdr-Dive)


This tells the validator to dive into a slice, array or map and validate that level of the slice, array or map with the validation tags that follow. Multidimensional nesting is also supported, each level you wish to dive will require another dive tag. dive has some sub-tags, 'keys' & 'endkeys', please see the Keys & EndKeys section just below. 


```
Usage: dive
```


Example #1 


```
[][]string with validation tag "gt=0,dive,len=1,dive,required"
// gt=0 will be applied to []
// len=1 will be applied to []string
// required will be applied to string
```


Example #2 


```
[][]string with validation tag "gt=0,dive,dive,required"
// gt=0 will be applied to []
// []string will be spared validation
// required will be applied to string
```


Keys & EndKeys 


These are to be used together directly after the dive tag and tells the validator that anything between 'keys' and 'endkeys' applies to the keys of a map and not the values; think of it like the 'dive' tag, but for map keys instead of values. Multidimensional nesting is also supported, each level you wish to validate will require another 'keys' and 'endkeys' tag. These tags are only valid for maps. 


```
Usage: dive,keys,othertagvalidation(s),endkeys,valuevalidationtags
```


Example #1 


```
map[string]string with validation tag "gt=0,dive,keys,eq=1|eq=2,endkeys,required"
// gt=0 will be applied to the map itself
// eq=1|eq=2 will be applied to the map keys
// required will be applied to map values
```


Example #2 


```
map[[2]string]string with validation tag "gt=0,dive,keys,dive,eq=1|eq=2,endkeys,required"
// gt=0 will be applied to the map itself
// eq=1|eq=2 will be applied to each array element in the map keys
// required will be applied to map values
```


#### Required [¶](#hdr-Required)


This validates that the value is not the data types default zero value. For numbers ensures value is not zero. For strings ensures value is not "". For booleans ensures value is not false. For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value when using WithRequiredStructEnabled. 


```
Usage: required
```


#### Required If [¶](#hdr-Required_If)


The field under validation must be present and not empty only if all the other specified fields are equal to the value following the specified field. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_if
```


Examples: 


```
// require the field if the Field1 is equal to the parameter given:
Usage: required_if=Field1 foobar

// require the field if the Field1 and Field2 is equal to the value respectively:
Usage: required_if=Field1 foo Field2 bar
```


#### Required Unless [¶](#hdr-Required_Unless)


The field under validation must be present and not empty unless all the other specified fields are equal to the value following the specified field. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_unless
```


Examples: 


```
// require the field unless the Field1 is equal to the parameter given:
Usage: required_unless=Field1 foobar

// require the field unless the Field1 and Field2 is equal to the value respectively:
Usage: required_unless=Field1 foo Field2 bar
```


#### Required With [¶](#hdr-Required_With)


The field under validation must be present and not empty only if any of the other specified fields are present. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_with
```


Examples: 


```
// require the field if the Field1 is present:
Usage: required_with=Field1

// require the field if the Field1 or Field2 is present:
Usage: required_with=Field1 Field2
```


#### Required With All [¶](#hdr-Required_With_All)


The field under validation must be present and not empty only if all of the other specified fields are present. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_with_all
```


Example: 


```
// require the field if the Field1 and Field2 is present:
Usage: required_with_all=Field1 Field2
```


#### Required Without [¶](#hdr-Required_Without)


The field under validation must be present and not empty only when any of the other specified fields are not present. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_without
```


Examples: 


```
// require the field if the Field1 is not present:
Usage: required_without=Field1

// require the field if the Field1 or Field2 is not present:
Usage: required_without=Field1 Field2
```


#### Required Without All [¶](#hdr-Required_Without_All)


The field under validation must be present and not empty only when all of the other specified fields are not present. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: required_without_all
```


Example: 


```
// require the field if the Field1 and Field2 is not present:
Usage: required_without_all=Field1 Field2
```


#### Excluded If [¶](#hdr-Excluded_If)


The field under validation must not be present or not empty only if all the other specified fields are equal to the value following the specified field. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: excluded_if
```


Examples: 


```
// exclude the field if the Field1 is equal to the parameter given:
Usage: excluded_if=Field1 foobar

// exclude the field if the Field1 and Field2 is equal to the value respectively:
Usage: excluded_if=Field1 foo Field2 bar
```


#### Excluded Unless [¶](#hdr-Excluded_Unless)


The field under validation must not be present or empty unless all the other specified fields are equal to the value following the specified field. For strings ensures value is not "". For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil. For structs ensures value is not the zero value. 


```
Usage: excluded_unless
```


Examples: 


```
// exclude the field unless the Field1 is equal to the parameter given:
Usage: excluded_unless=Field1 foobar

// exclude the field unless the Field1 and Field2 is equal to the value respectively:
Usage: excluded_unless=Field1 foo Field2 bar
```


#### Is Default [¶](#hdr-Is_Default)


This validates that the value is the default value and is almost the opposite of required. 


```
Usage: isdefault
```


#### Length [¶](#hdr-Length)


For numbers, length will ensure that the value is equal to the parameter given. For strings, it checks that the string length is exactly that number of characters. For slices, arrays, and maps, validates the number of items. 


Example #1 


```
Usage: len=10
```


Example #2 (time.Duration) 


For time.Duration, len will ensure that the value is equal to the duration given in the parameter. 


```
Usage: len=1h30m
```


#### Maximum [¶](#hdr-Maximum)


For numbers, max will ensure that the value is less than or equal to the parameter given. For strings, it checks that the string length is at most that number of characters. For slices, arrays, and maps, validates the number of items. 


Example #1 


```
Usage: max=10
```


Example #2 (time.Duration) 


For time.Duration, max will ensure that the value is less than or equal to the duration given in the parameter. 


```
Usage: max=1h30m
```


#### Minimum [¶](#hdr-Minimum)


For numbers, min will ensure that the value is greater or equal to the parameter given. For strings, it checks that the string length is at least that number of characters. For slices, arrays, and maps, validates the number of items. 


Example #1 


```
Usage: min=10
```


Example #2 (time.Duration) 


For time.Duration, min will ensure that the value is greater than or equal to the duration given in the parameter. 


```
Usage: min=1h30m
```


#### Equals [¶](#hdr-Equals)


For strings & numbers, eq will ensure that the value is equal to the parameter given. For slices, arrays, and maps, validates the number of items. 


Example #1 


```
Usage: eq=10
```


Example #2 (time.Duration) 


For time.Duration, eq will ensure that the value is equal to the duration given in the parameter. 


```
Usage: eq=1h30m
```


#### Not Equal [¶](#hdr-Not_Equal)


For strings & numbers, ne will ensure that the value is not equal to the parameter given. For slices, arrays, and maps, validates the number of items. 


Example #1 


```
Usage: ne=10
```


Example #2 (time.Duration) 


For time.Duration, ne will ensure that the value is not equal to the duration given in the parameter. 


```
Usage: ne=1h30m
```


#### One Of [¶](#hdr-One_Of)


For strings, ints, and uints, oneof will ensure that the value is one of the values in the parameter. The parameter should be a list of values separated by whitespace. Values may be strings or numbers. To match strings with spaces in them, include the target string between single quotes. Kind of like an 'enum'. 


```
Usage: oneof=red green
       oneof='red green' 'blue yellow'
       oneof=5 7 9
```


#### One Of Case Insensitive [¶](#hdr-One_Of_Case_Insensitive)


Works the same as oneof but is case insensitive and therefore only accepts strings. 


```
Usage: oneofci=red green
       oneofci='red green' 'blue yellow'
```


#### Greater Than [¶](#hdr-Greater_Than)


For numbers, this will ensure that the value is greater than the parameter given. For strings, it checks that the string length is greater than that number of characters. For slices, arrays and maps it validates the number of items. 


Example #1 


```
Usage: gt=10
```


Example #2 (time.Time) 


For time.Time ensures the time value is greater than time.Now.UTC(). 


```
Usage: gt
```


Example #3 (time.Duration) 


For time.Duration, gt will ensure that the value is greater than the duration given in the parameter. 


```
Usage: gt=1h30m
```


#### Greater Than or Equal [¶](#hdr-Greater_Than_or_Equal)


Same as 'min' above. Kept both to make terminology with 'len' easier. 


Example #1 


```
Usage: gte=10
```


Example #2 (time.Time) 


For time.Time ensures the time value is greater than or equal to time.Now.UTC(). 


```
Usage: gte
```


Example #3 (time.Duration) 


For time.Duration, gte will ensure that the value is greater than or equal to the duration given in the parameter. 


```
Usage: gte=1h30m
```


#### Less Than [¶](#hdr-Less_Than)


For numbers, this will ensure that the value is less than the parameter given. For strings, it checks that the string length is less than that number of characters. For slices, arrays, and maps it validates the number of items. 


Example #1 


```
Usage: lt=10
```


Example #2 (time.Time) 


For time.Time ensures the time value is less than time.Now.UTC(). 


```
Usage: lt
```


Example #3 (time.Duration) 


For time.Duration, lt will ensure that the value is less than the duration given in the parameter. 


```
Usage: lt=1h30m
```


#### Less Than or Equal [¶](#hdr-Less_Than_or_Equal)


Same as 'max' above. Kept both to make terminology with 'len' easier. 


Example #1 


```
Usage: lte=10
```


Example #2 (time.Time) 


For time.Time ensures the time value is less than or equal to time.Now.UTC(). 


```
Usage: lte
```


Example #3 (time.Duration) 


For time.Duration, lte will ensure that the value is less than or equal to the duration given in the parameter. 


```
Usage: lte=1h30m
```


#### Field Equals Another Field [¶](#hdr-Field_Equals_Another_Field)


This will validate the field value against another fields value either within a struct or passed in field. 


Example #1: 


```
// Validation on Password field using:
Usage: eqfield=ConfirmPassword
```


Example #2: 


```
// Validating by field:
validate.VarWithValue(password, confirmpassword, "eqfield")
```


Field Equals Another Field (relative) 


This does the same as eqfield except that it validates the field provided relative to the top level struct. 


```
Usage: eqcsfield=InnerStructField.Field)
```


#### Field Does Not Equal Another Field [¶](#hdr-Field_Does_Not_Equal_Another_Field)


This will validate the field value against another fields value either within a struct or passed in field. 


Examples: 


```
// Confirm two colors are not the same:
//
// Validation on Color field:
Usage: nefield=Color2

// Validating by field:
validate.VarWithValue(color1, color2, "nefield")
```


Field Does Not Equal Another Field (relative) 


This does the same as nefield except that it validates the field provided relative to the top level struct. 


```
Usage: necsfield=InnerStructField.Field
```


#### Field Greater Than Another Field [¶](#hdr-Field_Greater_Than_Another_Field)


Only valid for Numbers, time.Duration and time.Time types, this will validate the field value against another fields value either within a struct or passed in field. usage examples are for validation of a Start and End date: 


Example #1: 


```
// Validation on End field using:
validate.Struct Usage(gtfield=Start)
```


Example #2: 


```
// Validating by field:
validate.VarWithValue(start, end, "gtfield")
```


#### Field Greater Than Another Relative Field [¶](#hdr-Field_Greater_Than_Another_Relative_Field)


This does the same as gtfield except that it validates the field provided relative to the top level struct. 


```
Usage: gtcsfield=InnerStructField.Field
```


#### Field Greater Than or Equal To Another Field [¶](#hdr-Field_Greater_Than_or_Equal_To_Another_Field)


Only valid for Numbers, time.Duration and time.Time types, this will validate the field value against another fields value either within a struct or passed in field. usage examples are for validation of a Start and End date: 


Example #1: 


```
// Validation on End field using:
validate.Struct Usage(gtefield=Start)
```


Example #2: 


```
// Validating by field:
validate.VarWithValue(start, end, "gtefield")
```


#### Field Greater Than or Equal To Another Relative Field [¶](#hdr-Field_Greater_Than_or_Equal_To_Another_Relative_Field)


This does the same as gtefield except that it validates the field provided relative to the top level struct. 


```
Usage: gtecsfield=InnerStructField.Field
```


#### Less Than Another Field [¶](#hdr-Less_Than_Another_Field)


Only valid for Numbers, time.Duration and time.Time types, this will validate the field value against another fields value either within a struct or passed in field. usage examples are for validation of a Start and End date: 


Example #1: 


```
// Validation on End field using:
validate.Struct Usage(ltfield=Start)
```


Example #2: 


```
// Validating by field:
validate.VarWithValue(start, end, "ltfield")
```


#### Less Than Another Relative Field [¶](#hdr-Less_Than_Another_Relative_Field)


This does the same as ltfield except that it validates the field provided relative to the top level struct. 


```
Usage: ltcsfield=InnerStructField.Field
```


#### Less Than or Equal To Another Field [¶](#hdr-Less_Than_or_Equal_To_Another_Field)


Only valid for Numbers, time.Duration and time.Time types, this will validate the field value against another fields value either within a struct or passed in field. usage examples are for validation of a Start and End date: 


Example #1: 


```
// Validation on End field using:
validate.Struct Usage(ltefield=Start)
```


Example #2: 


```
// Validating by field:
validate.VarWithValue(start, end, "ltefield")
```


#### Less Than or Equal To Another Relative Field [¶](#hdr-Less_Than_or_Equal_To_Another_Relative_Field)


This does the same as ltefield except that it validates the field provided relative to the top level struct. 


```
Usage: ltecsfield=InnerStructField.Field
```


#### Field Contains Another Field [¶](#hdr-Field_Contains_Another_Field)


This does the same as contains except for struct fields. It should only be used with string types. See the behavior of reflect.Value.String() for behavior on other types. 


```
Usage: containsfield=InnerStructField.Field
```


#### Field Excludes Another Field [¶](#hdr-Field_Excludes_Another_Field)


This does the same as excludes except for struct fields. It should only be used with string types. See the behavior of reflect.Value.String() for behavior on other types. 


```
Usage: excludesfield=InnerStructField.Field
```


#### Unique [¶](#hdr-Unique)


For arrays & slices, unique will ensure that there are no duplicates. For maps, unique will ensure that there are no duplicate values. For slices of struct, unique will ensure that there are no duplicate values in a field of the struct specified via a parameter. 


```
// For arrays, slices, and maps:
Usage: unique

// For slices of struct:
Usage: unique=field
```


#### ValidateFn [¶](#hdr-ValidateFn)


This validates that an object responds to a method that can return error or bool. By default it expects an interface `Validate() error` and check that the method does not return an error. Other methods can be specified using two signatures: If the method returns an error, it check if the return value is nil. If the method returns a boolean, it checks if the value is true. 


```
// to use the default method Validate() error
Usage: validateFn

// to use the custom method IsValid() bool (or error)
Usage: validateFn=IsValid
```


#### Alpha Only [¶](#hdr-Alpha_Only)


This validates that a string value contains ASCII alpha characters only 


```
Usage: alpha
```


#### Alphanumeric [¶](#hdr-Alphanumeric)


This validates that a string value contains ASCII alphanumeric characters only 


```
Usage: alphanum
```


#### Alpha Unicode [¶](#hdr-Alpha_Unicode)


This validates that a string value contains unicode alpha characters only 


```
Usage: alphaunicode
```


#### Alphanumeric Unicode [¶](#hdr-Alphanumeric_Unicode)


This validates that a string value contains unicode alphanumeric characters only 


```
Usage: alphanumunicode
```


#### Boolean [¶](#hdr-Boolean)


This validates that a string value can successfully be parsed into a boolean with strconv.ParseBool 


```
Usage: boolean
```


#### Number [¶](#hdr-Number)


This validates that a string value contains number values only. For integers or float it returns true. 


```
Usage: number
```


#### Numeric [¶](#hdr-Numeric)


This validates that a string value contains a basic numeric value. basic excludes exponents etc... for integers or float it returns true. 


```
Usage: numeric
```


#### Hexadecimal String [¶](#hdr-Hexadecimal_String)


This validates that a string value contains a valid hexadecimal. 


```
Usage: hexadecimal
```


#### Hexcolor String [¶](#hdr-Hexcolor_String)


This validates that a string value contains a valid hex color including hashtag (#) 


```
Usage: hexcolor
```


#### Lowercase String [¶](#hdr-Lowercase_String)


This validates that a string value contains only lowercase characters. An empty string is not a valid lowercase string. 


```
Usage: lowercase
```


#### Uppercase String [¶](#hdr-Uppercase_String)


This validates that a string value contains only uppercase characters. An empty string is not a valid uppercase string. 


```
Usage: uppercase
```


#### RGB String [¶](#hdr-RGB_String)


This validates that a string value contains a valid rgb color 


```
Usage: rgb
```


#### RGBA String [¶](#hdr-RGBA_String)


This validates that a string value contains a valid rgba color 


```
Usage: rgba
```


#### HSL String [¶](#hdr-HSL_String)


This validates that a string value contains a valid hsl color 


```
Usage: hsl
```


#### HSLA String [¶](#hdr-HSLA_String)


This validates that a string value contains a valid hsla color 


```
Usage: hsla
```


#### E.164 Phone Number String [¶](#hdr-E_164_Phone_Number_String)


This validates that a string value contains a valid E.164 Phone number [https://en.wikipedia.org/wiki/E.164](https://en.wikipedia.org/wiki/E.164) (ex. +1123456789) 


```
Usage: e164
```


#### E-mail String [¶](#hdr-E_mail_String)


This validates that a string value contains a valid email This may not conform to all possibilities of any rfc standard, but neither does any email provider accept all possibilities. 


```
Usage: email
```


#### JSON String [¶](#hdr-JSON_String)


This validates that a string value is valid JSON 


```
Usage: json
```


#### JWT String [¶](#hdr-JWT_String)


This validates that a string value is a valid JWT 


```
Usage: jwt
```


#### File [¶](#hdr-File)


This validates that a string value contains a valid file path and that the file exists on the machine. This is done using os.Stat, which is a platform independent function. 


```
Usage: file
```


#### Image path [¶](#hdr-Image_path)


This validates that a string value contains a valid file path and that the file exists on the machine and is an image. This is done using os.Stat and github.com/gabriel-vasile/mimetype 


```
Usage: image
```


#### File Path [¶](#hdr-File_Path)


This validates that a string value contains a valid file path but does not validate the existence of that file. This is done using os.Stat, which is a platform independent function. 


```
Usage: filepath
```


#### URL String [¶](#hdr-URL_String)


This validates that a string value contains a valid url This will accept any url the golang request uri accepts but must contain a schema for example http:// or rtmp:// 


```
Usage: url
```


#### URI String [¶](#hdr-URI_String)


This validates that a string value contains a valid uri This will accept any uri the golang request uri accepts 


```
Usage: uri
```


#### Urn [RFC 2141](https://rfc-editor.org/rfc/rfc2141.html) String [RFC 2141](#hdr-Urn_RFC_2141_String) String" aria-label="Go to Urn [RFC 2141](https://rfc-editor.org/rfc/rfc2141.html) String">¶ 


This validates that a string value contains a valid URN according to the [RFC 2141](https://rfc-editor.org/rfc/rfc2141.html) spec. 


```
Usage: urn_rfc2141
```


#### Base32 String [¶](#hdr-Base32_String)


This validates that a string value contains a valid bas324 value. Although an empty string is valid base32 this will report an empty string as an error, if you wish to accept an empty string as valid you can use this with the omitempty tag. 


```
Usage: base32
```


#### Base64 String [¶](#hdr-Base64_String)


This validates that a string value contains a valid base64 value. Although an empty string is valid base64 this will report an empty string as an error, if you wish to accept an empty string as valid you can use this with the omitempty tag. 


```
Usage: base64
```


#### Base64URL String [¶](#hdr-Base64URL_String)


This validates that a string value contains a valid base64 URL safe value according the RFC4648 spec. Although an empty string is a valid base64 URL safe value, this will report an empty string as an error, if you wish to accept an empty string as valid you can use this with the omitempty tag. 


```
Usage: base64url
```


#### Base64RawURL String [¶](#hdr-Base64RawURL_String)


This validates that a string value contains a valid base64 URL safe value, but without = padding, according the RFC4648 spec, section 3.2. Although an empty string is a valid base64 URL safe value, this will report an empty string as an error, if you wish to accept an empty string as valid you can use this with the omitempty tag. 


```
Usage: base64rawurl
```


#### Bitcoin Address [¶](#hdr-Bitcoin_Address)


This validates that a string value contains a valid bitcoin address. The format of the string is checked to ensure it matches one of the three formats P2PKH, P2SH and performs checksum validation. 


```
Usage: btc_addr
```


Bitcoin Bech32 Address (segwit) 


This validates that a string value contains a valid bitcoin Bech32 address as defined by bip-0173 ( [https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki)) Special thanks to Pieter Wuille for providing reference implementations. 


```
Usage: btc_addr_bech32
```


#### Ethereum Address [¶](#hdr-Ethereum_Address)


This validates that a string value contains a valid ethereum address. The format of the string is checked to ensure it matches the standard Ethereum address format. 


```
Usage: eth_addr
```


#### Contains [¶](#hdr-Contains)


This validates that a string value contains the substring value. 


```
Usage: contains=@
```


#### Contains Any [¶](#hdr-Contains_Any)


This validates that a string value contains any Unicode code points in the substring value. 


```
Usage: containsany=!@#?
```


#### Contains Rune [¶](#hdr-Contains_Rune)


This validates that a string value contains the supplied rune value. 


```
Usage: containsrune=@
```


#### Excludes [¶](#hdr-Excludes)


This validates that a string value does not contain the substring value. 


```
Usage: excludes=@
```


#### Excludes All [¶](#hdr-Excludes_All)


This validates that a string value does not contain any Unicode code points in the substring value. 


```
Usage: excludesall=!@#?
```


#### Excludes Rune [¶](#hdr-Excludes_Rune)


This validates that a string value does not contain the supplied rune value. 


```
Usage: excludesrune=@
```


#### Starts With [¶](#hdr-Starts_With)


This validates that a string value starts with the supplied string value 


```
Usage: startswith=hello
```


#### Ends With [¶](#hdr-Ends_With)


This validates that a string value ends with the supplied string value 


```
Usage: endswith=goodbye
```


#### Does Not Start With [¶](#hdr-Does_Not_Start_With)


This validates that a string value does not start with the supplied string value 


```
Usage: startsnotwith=hello
```


#### Does Not End With [¶](#hdr-Does_Not_End_With)


This validates that a string value does not end with the supplied string value 


```
Usage: endsnotwith=goodbye
```


#### International Standard Book Number [¶](#hdr-International_Standard_Book_Number)


This validates that a string value contains a valid isbn10 or isbn13 value. 


```
Usage: isbn
```


#### International Standard Book Number 10 [¶](#hdr-International_Standard_Book_Number_10)


This validates that a string value contains a valid isbn10 value. 


```
Usage: isbn10
```


#### International Standard Book Number 13 [¶](#hdr-International_Standard_Book_Number_13)


This validates that a string value contains a valid isbn13 value. 


```
Usage: isbn13
```


#### Universally Unique Identifier UUID [¶](#hdr-Universally_Unique_Identifier_UUID)


This validates that a string value contains a valid UUID. Uppercase UUID values will not pass - use `uuid_rfc4122` instead. 


```
Usage: uuid
```


#### Universally Unique Identifier UUID v3 [¶](#hdr-Universally_Unique_Identifier_UUID_v3)


This validates that a string value contains a valid version 3 UUID. Uppercase UUID values will not pass - use `uuid3_rfc4122` instead. 


```
Usage: uuid3
```


#### Universally Unique Identifier UUID v4 [¶](#hdr-Universally_Unique_Identifier_UUID_v4)


This validates that a string value contains a valid version 4 UUID. Uppercase UUID values will not pass - use `uuid4_rfc4122` instead. 


```
Usage: uuid4
```


#### Universally Unique Identifier UUID v5 [¶](#hdr-Universally_Unique_Identifier_UUID_v5)


This validates that a string value contains a valid version 5 UUID. Uppercase UUID values will not pass - use `uuid5_rfc4122` instead. 


```
Usage: uuid5
```


#### Universally Unique Lexicographically Sortable Identifier ULID [¶](#hdr-Universally_Unique_Lexicographically_Sortable_Identifier_ULID)


This validates that a string value contains a valid ULID value. 


```
Usage: ulid
```


#### ASCII [¶](#hdr-ASCII)


This validates that a string value contains only ASCII characters. NOTE: if the string is blank, this validates as true. 


```
Usage: ascii
```


#### Printable ASCII [¶](#hdr-Printable_ASCII)


This validates that a string value contains only printable ASCII characters. NOTE: if the string is blank, this validates as true. 


```
Usage: printascii
```


#### Multi-Byte Characters [¶](#hdr-Multi_Byte_Characters)


This validates that a string value contains one or more multibyte characters. NOTE: if the string is blank, this validates as true. 


```
Usage: multibyte
```


#### Data URL [¶](#hdr-Data_URL)


This validates that a string value contains a valid DataURI. NOTE: this will also validate that the data portion is valid base64 


```
Usage: datauri
```


#### Latitude [¶](#hdr-Latitude)


This validates that a string value contains a valid latitude. 


```
Usage: latitude
```


#### Longitude [¶](#hdr-Longitude)


This validates that a string value contains a valid longitude. 


```
Usage: longitude
```


#### Employeer Identification Number EIN [¶](#hdr-Employeer_Identification_Number_EIN)


This validates that a string value contains a valid U.S. Employer Identification Number. 


```
Usage: ein
```


#### Social Security Number SSN [¶](#hdr-Social_Security_Number_SSN)


This validates that a string value contains a valid U.S. Social Security Number. 


```
Usage: ssn
```


#### Internet Protocol Address IP [¶](#hdr-Internet_Protocol_Address_IP)


This validates that a string value contains a valid IP Address. 


```
Usage: ip
```


#### Internet Protocol Address IPv4 [¶](#hdr-Internet_Protocol_Address_IPv4)


This validates that a string value contains a valid v4 IP Address. 


```
Usage: ipv4
```


#### Internet Protocol Address IPv6 [¶](#hdr-Internet_Protocol_Address_IPv6)


This validates that a string value contains a valid v6 IP Address. 


```
Usage: ipv6
```


#### Classless Inter-Domain Routing CIDR [¶](#hdr-Classless_Inter_Domain_Routing_CIDR)


This validates that a string value contains a valid CIDR Address. 


```
Usage: cidr
```


#### Classless Inter-Domain Routing CIDRv4 [¶](#hdr-Classless_Inter_Domain_Routing_CIDRv4)


This validates that a string value contains a valid v4 CIDR Address. 


```
Usage: cidrv4
```


#### Classless Inter-Domain Routing CIDRv6 [¶](#hdr-Classless_Inter_Domain_Routing_CIDRv6)


This validates that a string value contains a valid v6 CIDR Address. 


```
Usage: cidrv6
```


#### Transmission Control Protocol Address TCP [¶](#hdr-Transmission_Control_Protocol_Address_TCP)


This validates that a string value contains a valid resolvable TCP Address. 


```
Usage: tcp_addr
```


#### Transmission Control Protocol Address TCPv4 [¶](#hdr-Transmission_Control_Protocol_Address_TCPv4)


This validates that a string value contains a valid resolvable v4 TCP Address. 


```
Usage: tcp4_addr
```


#### Transmission Control Protocol Address TCPv6 [¶](#hdr-Transmission_Control_Protocol_Address_TCPv6)


This validates that a string value contains a valid resolvable v6 TCP Address. 


```
Usage: tcp6_addr
```


#### User Datagram Protocol Address UDP [¶](#hdr-User_Datagram_Protocol_Address_UDP)


This validates that a string value contains a valid resolvable UDP Address. 


```
Usage: udp_addr
```


#### User Datagram Protocol Address UDPv4 [¶](#hdr-User_Datagram_Protocol_Address_UDPv4)


This validates that a string value contains a valid resolvable v4 UDP Address. 


```
Usage: udp4_addr
```


#### User Datagram Protocol Address UDPv6 [¶](#hdr-User_Datagram_Protocol_Address_UDPv6)


This validates that a string value contains a valid resolvable v6 UDP Address. 


```
Usage: udp6_addr
```


#### Internet Protocol Address IP [¶](#hdr-Internet_Protocol_Address_IP-1)


This validates that a string value contains a valid resolvable IP Address. 


```
Usage: ip_addr
```


#### Internet Protocol Address IPv4 [¶](#hdr-Internet_Protocol_Address_IPv4-1)


This validates that a string value contains a valid resolvable v4 IP Address. 


```
Usage: ip4_addr
```


#### Internet Protocol Address IPv6 [¶](#hdr-Internet_Protocol_Address_IPv6-1)


This validates that a string value contains a valid resolvable v6 IP Address. 


```
Usage: ip6_addr
```


#### Unix domain socket end point Address [¶](#hdr-Unix_domain_socket_end_point_Address)


This validates that a string value contains a valid Unix Address. 


```
Usage: unix_addr
```


#### Media Access Control Address MAC [¶](#hdr-Media_Access_Control_Address_MAC)


This validates that a string value contains a valid MAC Address. 


```
Usage: mac
```


Note: See Go's ParseMAC for accepted formats and types: 


```
http://golang.org/src/net/mac.go?s=866:918#L29
```


#### Hostname [RFC 952](https://rfc-editor.org/rfc/rfc952.html)
                [RFC 952](#hdr-Hostname_RFC_952)" aria-label="Go to Hostname [RFC 952](https://rfc-editor.org/rfc/rfc952.html)">¶


This validates that a string value is a valid Hostname according to [RFC 952](https://rfc-editor.org/rfc/rfc952.html)
                [https://tools.ietf.org/html/rfc952](https://tools.ietf.org/html/rfc952)


```
Usage: hostname
```


#### Hostname [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html)
                [RFC 1123](#hdr-Hostname_RFC_1123)" aria-label="Go to Hostname [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html)">¶


This validates that a string value is a valid Hostname according to [RFC 1123](https://rfc-editor.org/rfc/rfc1123.html)
                [https://tools.ietf.org/html/rfc1123](https://tools.ietf.org/html/rfc1123)


```
Usage: hostname_rfc1123 or if you want to continue to use 'hostname' in your tags, create an alias.
```


Full Qualified Domain Name (FQDN) 


This validates that a string value contains a valid FQDN. 


```
Usage: fqdn
```


#### HTML Tags [¶](#hdr-HTML_Tags)


This validates that a string value appears to be an HTML element tag including those described at [https://developer.mozilla.org/en-US/docs/Web/HTML/Element](https://developer.mozilla.org/en-US/docs/Web/HTML/Element)


```
Usage: html
```


#### HTML Encoded [¶](#hdr-HTML_Encoded)


This validates that a string value is a proper character reference in decimal or hexadecimal format 


```
Usage: html_encoded
```


#### URL Encoded [¶](#hdr-URL_Encoded)


This validates that a string value is percent-encoded (URL encoded) according to [https://tools.ietf.org/html/rfc3986#section-2.1](https://tools.ietf.org/html/rfc3986#section-2.1)


```
Usage: url_encoded
```


#### Directory [¶](#hdr-Directory)


This validates that a string value contains a valid directory and that it exists on the machine. This is done using os.Stat, which is a platform independent function. 


```
Usage: dir
```


#### Directory Path [¶](#hdr-Directory_Path)


This validates that a string value contains a valid directory but does not validate the existence of that directory. This is done using os.Stat, which is a platform independent function. It is safest to suffix the string with os.PathSeparator if the directory may not exist at the time of validation. 


```
Usage: dirpath
```


#### HostPort [¶](#hdr-HostPort)


This validates that a string value contains a valid DNS hostname and port that can be used to validate fields typically passed to sockets and connections. 


```
Usage: hostname_port
```


#### Datetime [¶](#hdr-Datetime)


This validates that a string value is a valid datetime based on the supplied datetime format. Supplied format must match the official Go time format layout as documented in [https://golang.org/pkg/time/](https://golang.org/pkg/time/)


```
Usage: datetime=2006-01-02
```


#### Iso3166-1 alpha-2 [¶](#hdr-Iso3166_1_alpha_2)


This validates that a string value is a valid country code based on iso3166-1 alpha-2 standard. see: [https://www.iso.org/iso-3166-country-codes.html](https://www.iso.org/iso-3166-country-codes.html)


```
Usage: iso3166_1_alpha2
```


#### Iso3166-1 alpha-3 [¶](#hdr-Iso3166_1_alpha_3)


This validates that a string value is a valid country code based on iso3166-1 alpha-3 standard. see: [https://www.iso.org/iso-3166-country-codes.html](https://www.iso.org/iso-3166-country-codes.html)


```
Usage: iso3166_1_alpha3
```


#### Iso3166-1 alpha-numeric [¶](#hdr-Iso3166_1_alpha_numeric)


This validates that a string value is a valid country code based on iso3166-1 alpha-numeric standard. see: [https://www.iso.org/iso-3166-country-codes.html](https://www.iso.org/iso-3166-country-codes.html)


```
Usage: iso3166_1_alpha3
```


#### BCP 47 Language Tag [¶](#hdr-BCP_47_Language_Tag)


This validates that a string value is a valid BCP 47 language tag, as parsed by language.Parse. More information on [https://pkg.go.dev/golang.org/x/text/language](https://pkg.go.dev/golang.org/x/text/language)


```
Usage: bcp47_language_tag
```


BIC (SWIFT code) 


This validates that a string value is a valid Business Identifier Code (SWIFT code), defined in ISO 9362. More information on [https://www.iso.org/standard/60390.html](https://www.iso.org/standard/60390.html)


```
Usage: bic
```


#### 
                [RFC 1035](https://rfc-editor.org/rfc/rfc1035.html) label [RFC 1035](#hdr-RFC_1035_label) label" aria-label="Go to [RFC 1035](https://rfc-editor.org/rfc/rfc1035.html) label">¶


This validates that a string value is a valid dns [RFC 1035](https://rfc-editor.org/rfc/rfc1035.html) label, defined in [RFC 1035](https://rfc-editor.org/rfc/rfc1035.html). More information on [https://datatracker.ietf.org/doc/html/rfc1035](https://datatracker.ietf.org/doc/html/rfc1035)


```
Usage: dns_rfc1035_label
```


#### TimeZone [¶](#hdr-TimeZone)


This validates that a string value is a valid time zone based on the time zone database present on the system. Although empty value and Local value are allowed by time.LoadLocation golang function, they are not allowed by this validator. More information on [https://golang.org/pkg/time/#LoadLocation](https://golang.org/pkg/time/#LoadLocation)


```
Usage: timezone
```


#### Semantic Version [¶](#hdr-Semantic_Version)


This validates that a string value is a valid semver version, defined in Semantic Versioning 2.0.0. More information on [https://semver.org/](https://semver.org/)


```
Usage: semver
```


#### CVE Identifier [¶](#hdr-CVE_Identifier)


This validates that a string value is a valid cve id, defined in cve mitre. More information on [https://cve.mitre.org/](https://cve.mitre.org/)


```
Usage: cve
```


#### Credit Card [¶](#hdr-Credit_Card)


This validates that a string value contains a valid credit card number using Luhn algorithm. 


```
Usage: credit_card
```


#### Luhn Checksum [¶](#hdr-Luhn_Checksum)


```
Usage: luhn_checksum
```


This validates that a string or (u)int value contains a valid checksum using the Luhn algorithm. 


#### MongoDB [¶](#hdr-MongoDB)


This validates that a string is a valid 24 character hexadecimal string or valid connection string. 


```
Usage: mongodb
       mongodb_connection_string
```


Example: 


```
type Test struct {
  ObjectIdField         string `validate:"mongodb"`
  ConnectionStringField string `validate:"mongodb_connection_string"`
}
```


#### Cron [¶](#hdr-Cron)


This validates that a string value contains a valid cron expression. 


```
Usage: cron
```


#### SpiceDb ObjectID/Permission/Object Type [¶](#hdr-SpiceDb_ObjectID_Permission_Object_Type)


This validates that a string is valid for use with SpiceDb for the indicated purpose. If no purpose is given, a purpose of 'id' is assumed. 


```
Usage: spicedb=id|permission|type
```


#### Alias Validators and Tags [¶](#hdr-Alias_Validators_and_Tags)


Alias Validators and Tags NOTE: When returning an error, the tag returned in "FieldError" will be the alias tag unless the dive tag is part of the alias. Everything after the dive tag is not reported as the alias tag. Also, the "ActualTag" in the before case will be the actual tag within the alias that failed. 


Here is a list of the current built in alias tags: 


```
"iscolor"
  alias is "hexcolor|rgb|rgba|hsl|hsla" (Usage: iscolor)
"country_code"
  alias is "iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric" (Usage: country_code)
```


Validator notes: 


```
regex
  a regex validator won't be added because commas and = signs can be part
  of a regex which conflict with the validation definitions. Although
  workarounds can be made, they take away from using pure regex's.
  Furthermore it's quick and dirty but the regex's become harder to
  maintain and are not reusable, so it's as much a programming philosophy
  as anything.

  In place of this new validator functions should be created; a regex can
  be used within the validator function and even be precompiled for better
  efficiency within regexes.go.

  And the best reason, you can submit a pull request and we can keep on
  adding to the validation library of this package!
```


#### Non standard validators [¶](#hdr-Non_standard_validators)


A collection of validation rules that are frequently needed but are more complex than the ones found in the baked in validators. A non standard validator must be registered manually like you would with your own custom validation functions. 


Example of registration and use: 


```
type Test struct {
  TestField string `validate:"yourtag"`
}

t := &Test{
  TestField: "Test"
}

validate := validator.New()
validate.RegisterValidation("yourtag", validators.NotBlank)
```


Here is a list of the current non standard validators: 


```
NotBlank
  This validates that the value is not blank or with length zero.
  For strings ensures they do not contain only spaces. For channels, maps, slices and arrays
  ensures they don't have zero length. For others, a non empty value is required.

  Usage: notblank
```


#### Panics [¶](#hdr-Panics)


This package panics when bad input is provided, this is by design, bad code like that should not make it to production. 


```
type Test struct {
  TestField string `validate:"nonexistantfunction=1"`
}

t := &Test{
  TestField: "Test"
}

validate.Struct(t) // this will panic
```


### Index [¶](#pkg-index)


                - 
                  [type CustomTypeFunc](#CustomTypeFunc)


                - 
                  [type FieldError](#FieldError)


                - 
                  [type FieldLevel](#FieldLevel)


                - 
                  [type FilterFunc](#FilterFunc)


                - 
                  [type Func](#Func)


                - 
                  [type FuncCtx](#FuncCtx)


                - 
                  [type InvalidValidationError](#InvalidValidationError)


                - 


                      - 
                      [func (e *InvalidValidationError) Error() string](#InvalidValidationError.Error)


                - 
                  [type Option](#Option)


                - 


                      - 
                      [func WithPrivateFieldValidation() Option](#WithPrivateFieldValidation)


                      - 
                      [func WithRequiredStructEnabled() Option](#WithRequiredStructEnabled)


                - 
                  [type RegisterTranslationsFunc](#RegisterTranslationsFunc)


                - 
                  [type StructLevel](#StructLevel)


                - 
                  [type StructLevelFunc](#StructLevelFunc)


                - 
                  [type StructLevelFuncCtx](#StructLevelFuncCtx)


                - 
                  [type TagNameFunc](#TagNameFunc)


                - 
                  [type TranslationFunc](#TranslationFunc)


                - 
                  [type Validate](#Validate)


                - 


                      - 
                      [func New(options ...Option) *Validate](#New)


                - 


                      - 
                      [func (v *Validate) RegisterAlias(alias, tags string)](#Validate.RegisterAlias)


                      - 
                      [func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{})](#Validate.RegisterCustomTypeFunc)


                      - 
                      [func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...interface{})](#Validate.RegisterStructValidation)


                      - 
                      [func (v *Validate) RegisterStructValidationCtx(fn StructLevelFuncCtx, types ...interface{})](#Validate.RegisterStructValidationCtx)


                      - 
                      [func (v *Validate) RegisterStructValidationMapRules(rules map[string]string, types ...interface{})](#Validate.RegisterStructValidationMapRules)


                      - 
                      [func (v *Validate) RegisterTagNameFunc(fn TagNameFunc)](#Validate.RegisterTagNameFunc)


                      - 
                      [func (v *Validate) RegisterTranslation(tag string, trans ut.Translator, registerFn RegisterTranslationsFunc, ...) (err error)](#Validate.RegisterTranslation)


                      - 
                      [func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error](#Validate.RegisterValidation)


                      - 
                      [func (v *Validate) RegisterValidationCtx(tag string, fn FuncCtx, callValidationEvenIfNull ...bool) error](#Validate.RegisterValidationCtx)


                      - 
                      [func (v *Validate) SetTagName(name string)](#Validate.SetTagName)


                      - 
                      [func (v *Validate) Struct(s interface{}) error](#Validate.Struct)


                      - 
                      [func (v *Validate) StructCtx(ctx context.Context, s interface{}) (err error)](#Validate.StructCtx)


                      - 
                      [func (v *Validate) StructExcept(s interface{}, fields ...string) error](#Validate.StructExcept)


                      - 
                      [func (v *Validate) StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error)](#Validate.StructExceptCtx)


                      - 
                      [func (v *Validate) StructFiltered(s interface{}, fn FilterFunc) error](#Validate.StructFiltered)


                      - 
                      [func (v *Validate) StructFilteredCtx(ctx context.Context, s interface{}, fn FilterFunc) (err error)](#Validate.StructFilteredCtx)


                      - 
                      [func (v *Validate) StructPartial(s interface{}, fields ...string) error](#Validate.StructPartial)


                      - 
                      [func (v *Validate) StructPartialCtx(ctx context.Context, s interface{}, fields ...string) (err error)](#Validate.StructPartialCtx)


                      - 
                      [func (v *Validate) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{}](#Validate.ValidateMap)


                      - 
                      [func (v Validate) ValidateMapCtx(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{}](#Validate.ValidateMapCtx)


                      - 
                      [func (v *Validate) Var(field interface{}, tag string) error](#Validate.Var)


                      - 
                      [func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag string) (err error)](#Validate.VarCtx)


                      - 
                      [func (v *Validate) VarWithValue(field interface{}, other interface{}, tag string) error](#Validate.VarWithValue)


                      - 
                      [func (v *Validate) VarWithValueCtx(ctx context.Context, field interface{}, other interface{}, tag string) (err error)](#Validate.VarWithValueCtx)


                - 
                  [type ValidationErrors](#ValidationErrors)


                - 


                      - 
                      [func (ve ValidationErrors) Error() string](#ValidationErrors.Error)


                      - 
                      [func (ve ValidationErrors) Translate(ut ut.Translator) ValidationErrorsTranslations](#ValidationErrors.Translate)


                - 
                  [type ValidationErrorsTranslations](#ValidationErrorsTranslations)


### Constants [¶](#pkg-constants)


This section is empty.


### Variables [¶](#pkg-variables)


This section is empty.


### Functions [¶](#pkg-functions)


This section is empty.


### Types [¶](#pkg-types)


#### 
                  type [CustomTypeFunc](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L72)
                    [¶](#CustomTypeFunc)


                  [](/reflect)[](/reflect#Value)
```
type CustomTypeFunc func(field 
                        reflect.
                        Value) interface{}

```


CustomTypeFunc allows for overriding or adding custom field type handler functions field = field value of the type to return a value to be validated example Valuer from sql drive see [https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29](https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29)


#### 
                  type [FieldError](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L76)
                    [¶](#FieldError)


                  [](/builtin#string)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/reflect)[](/reflect#Kind)[](/reflect)[](/reflect#Type)[](/github.com/go-playground/universal-translator)[](/github.com/go-playground/universal-translator#Translator)[](/builtin#string)[](/builtin#string)
```
type FieldError interface {


                          // Tag returns the validation tag that failed. if the

                        // validation was an alias, this will return the
                        // alias name and not the underlying tag that failed.
                        //
                        // eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
                        // will return "iscolor"
  Tag() 
                        string

                          // ActualTag returns the validation tag that failed, even if an

                        // alias the actual tag within the alias will be returned.
                        // If an 'or' validation fails the entire or will be returned.
                        //
                        // eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
                        // will return "hexcolor|rgb|rgba|hsl|hsla"
  ActualTag() 
                        string

                          // Namespace returns the namespace for the field error, with the tag

                        // name taking precedence over the field's actual name.
                        //
                        // eg. JSON name "User.fname"
                        //
                        // See StructNamespace() for a version that returns actual names.
                        //
                        // NOTE: this field can be blank when validating a single primitive field
                        // using validate.Field(...) as there is no way to extract it's name
  Namespace() 
                        string

                          // StructNamespace returns the namespace for the field error, with the field's

                        // actual name.
                        //
                        // eg. "User.FirstName" see Namespace for comparison
                        //
                        // NOTE: this field can be blank when validating a single primitive field
                        // using validate.Field(...) as there is no way to extract its name
  StructNamespace() 
                        string

                          // Field returns the field's name with the tag name taking precedence over the

                        // field's actual name.
                        //
                        // `RegisterTagNameFunc` must be registered to get tag value.
                        //
                        // eg. JSON name "fname"
                        // see StructField for comparison
  Field() 
                        string

                          // StructField returns the field's actual name from the struct, when able to determine.

                        //
                        // eg.  "FirstName"
                        // see Field for comparison
  StructField() 
                        string

                          // Value returns the actual field's value in case needed for creating the error

                        // message
  Value() interface{}


                          // Param returns the param value, in string form for comparison; this will also

                        // help with generating an error message
  Param() 
                        string

                          // Kind returns the Field's reflect Kind

                        //
                        // eg. time.Time's kind is a struct
  Kind() 
                        reflect.
                        Kind

                          // Type returns the Field's reflect Type

                        //
                        // eg. time.Time's type is time.Time
  Type() 
                        reflect.
                        Type

                          // Translate returns the FieldError's translated error

                        // from the provided 'ut.Translator' and registered 'TranslationFunc'
                        //
                        // NOTE: if no registered translator can be found it returns the same as
                        // calling fe.Error()
  Translate(ut 
                        ut.
                        Translator) 
                        string

                          // Error returns the FieldError's message
                         Error() 
                        string
}

```


FieldError contains all functions to get error details 


#### 
                  type [FieldLevel](https://github.com/go-playground/validator/blob/v10.27.0/field_level.go#L7)
                    [¶](#FieldLevel)


                  [](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/reflect)[](/reflect#Value)[](/builtin#string)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/builtin#bool)[](/reflect)[](/reflect#Value)[](/builtin#string)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/builtin#bool)
```
type FieldLevel interface {


                          // Top returns the top level struct, if any
                         Top() 
                        reflect.
                        Value

                          // Parent returns the current fields parent struct, if any or

                        // the comparison value if called 'VarWithValue'
  Parent() 
                        reflect.
                        Value

                          // Field returns current field for validation
                         Field() 
                        reflect.
                        Value

                          // FieldName returns the field's name with the tag

                        // name taking precedence over the fields actual name.
  FieldName() 
                        string

                          // StructFieldName returns the struct field's name
                         StructFieldName() 
                        string

                          // Param returns param for validation against current field
                         Param() 
                        string

                          // GetTag returns the current validations tag name
                         GetTag() 
                        string

                          // ExtractType gets the actual underlying type of field value.

                        // It will dive into pointers, customTypes and return you the
                        // underlying value and it's kind.
  ExtractType(field 
                        reflect.
                        Value) (value 
                        reflect.
                        Value, kind 
                        reflect.
                        Kind, nullable 
                        bool)


                          // GetStructFieldOK traverses the parent struct to retrieve a specific field denoted by the provided namespace

                        // in the param and returns the field, field kind and whether is was successful in retrieving
                        // the field at all.
                        //
                        // NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
                        // could not be retrieved because it didn't exist.
                        //
                        // Deprecated: Use GetStructFieldOK2() instead which also return if the value is nullable.
  GetStructFieldOK() (
                        reflect.
                        Value, 
                        reflect.
                        Kind, 
                        bool)


                          // GetStructFieldOKAdvanced is the same as GetStructFieldOK except that it accepts the parent struct to start looking for

                        // the field and namespace allowing more extensibility for validators.
                        //
                        // Deprecated: Use GetStructFieldOKAdvanced2() instead which also return if the value is nullable.
  GetStructFieldOKAdvanced(val 
                        reflect.
                        Value, namespace 
                        string) (
                        reflect.
                        Value, 
                        reflect.
                        Kind, 
                        bool)


                          // GetStructFieldOK2 traverses the parent struct to retrieve a specific field denoted by the provided namespace

                        // in the param and returns the field, field kind, if it's a nullable type and whether is was successful in retrieving
                        // the field at all.
                        //
                        // NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
                        // could not be retrieved because it didn't exist.
  GetStructFieldOK2() (
                        reflect.
                        Value, 
                        reflect.
                        Kind, 
                        bool, 
                        bool)


                          // GetStructFieldOKAdvanced2 is the same as GetStructFieldOK except that it accepts the parent struct to start looking for

                        // the field and namespace allowing more extensibility for validators.
  GetStructFieldOKAdvanced2(val 
                        reflect.
                        Value, namespace 
                        string) (
                        reflect.
                        Value, 
                        reflect.
                        Kind, 
                        bool, 
                        bool)
}

```


FieldLevel contains all the information and helper functions to validate a field 


#### 
                  type [FilterFunc](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L67)
                    [¶](#FilterFunc)


                  [](/builtin#byte)[](/builtin#bool)
```
type FilterFunc func(ns []
                        byte) 
                        bool

```


FilterFunc is the type used to filter fields using StructFiltered(...) function. returning true results in the field being filtered/skipped from validation 


#### 
                  type [Func](https://github.com/go-playground/validator/blob/v10.27.0/baked_in.go#L34)
                    [¶](#Func)


                  [](#FieldLevel)[](/builtin#bool)
```
type Func func(fl 
                        FieldLevel) 
                        bool

```


Func accepts a FieldLevel interface for all validation needs. The return value should be true when validation succeeds. 


#### 
                  type [FuncCtx](https://github.com/go-playground/validator/blob/v10.27.0/baked_in.go#L38)
                    [¶](#FuncCtx)


                  [](/context)[](/context#Context)[](#FieldLevel)[](/builtin#bool)
```
type FuncCtx func(ctx 
                        context.
                        Context, fl 
                        FieldLevel) 
                        bool

```


FuncCtx accepts a context.Context and FieldLevel interface for all validation needs. The return value should be true when validation succeeds. 


#### 
                  type [InvalidValidationError](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L21)
                    [¶](#InvalidValidationError)


                  [](/reflect)[](/reflect#Type)
```
type InvalidValidationError struct {

                         Type 
                          reflect.
                          Type
                        }

```


InvalidValidationError describes an invalid argument passed to `Struct`, `StructExcept`, StructPartial` or `Field` 


#### 
                    func (*InvalidValidationError) [Error](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L26)
                      [¶](#InvalidValidationError.Error)


                    [](#InvalidValidationError)[](/builtin#string)
```
func (e *
                          InvalidValidationError) Error() 
                          string

```


Error returns InvalidValidationError message 


#### 
                  type [Option](https://github.com/go-playground/validator/blob/v10.27.0/options.go#L4)
                    [¶](#Option)


                    added in
                    v10.15.2


                  [](#Validate)
```
type Option func(*
                        Validate)

```


Option represents a configurations option to be applied to validator during initialization. 


#### 
                    func [WithPrivateFieldValidation](https://github.com/go-playground/validator/blob/v10.27.0/options.go#L22)
                      [¶](#WithPrivateFieldValidation)


                      added in
                      v10.19.0


                    [](#Option)
```
func WithPrivateFieldValidation() 
                          Option

```


WithPrivateFieldValidation activates validation for unexported fields via the use of the `unsafe` package. 


By opting into this feature you are acknowledging that you are aware of the risks and accept any current or future consequences of using this feature. 


#### 
                    func [WithRequiredStructEnabled](https://github.com/go-playground/validator/blob/v10.27.0/options.go#L12)
                      [¶](#WithRequiredStructEnabled)


                      added in
                      v10.15.2


                    [](#Option)
```
func WithRequiredStructEnabled() 
                          Option

```


WithRequiredStructEnabled enables required tag on non-pointer structs to be applied instead of ignored. 


This was made opt-in behaviour in order to maintain backward compatibility with the behaviour previous to being able to apply struct level validations on struct fields directly. 


It is recommended you enabled this as it will be the default behaviour in v11+ 


#### 
                  type [RegisterTranslationsFunc](https://github.com/go-playground/validator/blob/v10.27.0/translations.go#L11)
                    [¶](#RegisterTranslationsFunc)


                  [](/github.com/go-playground/universal-translator)[](/github.com/go-playground/universal-translator#Translator)[](/builtin#error)
```
type RegisterTranslationsFunc func(ut 
                        ut.
                        Translator) 
                        error

```


RegisterTranslationsFunc allows for registering of translations for a 'ut.Translator' for use within the 'TranslationFunc' 


#### 
                  type [StructLevel](https://github.com/go-playground/validator/blob/v10.27.0/struct_level.go#L24)
                    [¶](#StructLevel)


                  [](#Validate)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Value)[](/reflect)[](/reflect#Kind)[](/builtin#bool)[](/builtin#string)[](/builtin#string)[](/builtin#string)[](#ValidationErrors)
```
type StructLevel interface {


                          // Validator returns the main validation object, in case one wants to call validations internally.

                        // this is so you don't have to use anonymous functions to get access to the validate
                        // instance.
  Validator() *
                        Validate

                          // Top returns the top level struct, if any
                         Top() 
                        reflect.
                        Value

                          // Parent returns the current fields parent struct, if any
                         Parent() 
                        reflect.
                        Value

                          // Current returns the current struct.
                         Current() 
                        reflect.
                        Value

                          // ExtractType gets the actual underlying type of field value.

                        // It will dive into pointers, customTypes and return you the
                        // underlying value and its kind.
  ExtractType(field 
                        reflect.
                        Value) (value 
                        reflect.
                        Value, kind 
                        reflect.
                        Kind, nullable 
                        bool)


                          // ReportError reports an error just by passing the field and tag information

                        //
                        // NOTES:
                        //
                        // fieldName and structFieldName get appended to the existing
                        // namespace that validator is on. e.g. pass 'FirstName' or
                        // 'Names[0]' depending on the nesting
                        //
                        // tag can be an existing validation tag or just something you make up
                        // and process on the flip side it's up to you.
  ReportError(field interface{}, fieldName, structFieldName 
                        string, tag, param 
                        string)


                          // ReportValidationErrors reports an error just by passing ValidationErrors

                        //
                        // NOTES:
                        //
                        // relativeNamespace and relativeActualNamespace get appended to the
                        // existing namespace that validator is on.
                        // e.g. pass 'User.FirstName' or 'Users[0].FirstName' depending
                        // on the nesting. most of the time they will be blank, unless you validate
                        // at a level lower the current field depth
  ReportValidationErrors(relativeNamespace, relativeActualNamespace 
                        string, errs 
                        ValidationErrors)
}

```


StructLevel contains all the information and helper functions to validate a struct 


#### 
                  type [StructLevelFunc](https://github.com/go-playground/validator/blob/v10.27.0/struct_level.go#L9)
                    [¶](#StructLevelFunc)


                  [](#StructLevel)
```
type StructLevelFunc func(sl 
                        StructLevel)

```


StructLevelFunc accepts all values needed for struct level validation 


#### 
                  type [StructLevelFuncCtx](https://github.com/go-playground/validator/blob/v10.27.0/struct_level.go#L13)
                    [¶](#StructLevelFuncCtx)


                  [](/context)[](/context#Context)[](#StructLevel)
```
type StructLevelFuncCtx func(ctx 
                        context.
                        Context, sl 
                        StructLevel)

```


StructLevelFuncCtx accepts all values needed for struct level validation but also allows passing of contextual validation information via context.Context. 


#### 
                  type [TagNameFunc](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L75)
                    [¶](#TagNameFunc)


                  [](/reflect)[](/reflect#StructField)[](/builtin#string)
```
type TagNameFunc func(field 
                        reflect.
                        StructField) 
                        string

```


TagNameFunc allows for adding of a custom tag name parser 


#### 
                  type [TranslationFunc](https://github.com/go-playground/validator/blob/v10.27.0/translations.go#L7)
                    [¶](#TranslationFunc)


                  [](/github.com/go-playground/universal-translator)[](/github.com/go-playground/universal-translator#Translator)[](#FieldError)[](/builtin#string)
```
type TranslationFunc func(ut 
                        ut.
                        Translator, fe 
                        FieldError) 
                        string

```


TranslationFunc is the function type used to register or override custom translations 


#### 
                  type [Validate](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L83)
                    [¶](#Validate)


```
type Validate struct {

                        // contains filtered or unexported fields
}

```


Validate contains the validator settings and cache 


#### 
                    func [New](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L106)
                      [¶](#New)


                    [](#Option)[](#Validate)
```
func New(options ...
                          Option) *
                          Validate

```


New returns a new instance of 'validate' with sane defaults. Validate is designed to be thread-safe and used as a singleton instance. It caches information about your struct and validations, in essence only parsing your validation tags once per struct type. Using multiple instances neglects the benefit of caching. 


#### 
                    func (*Validate) [RegisterAlias](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L239)
                      [¶](#Validate.RegisterAlias)


                    [](#Validate)[](/builtin#string)
```
func (v *
                          Validate) RegisterAlias(alias, tags 
                          string)

```


RegisterAlias registers a mapping of a single validation tag that defines a common or complex set of validation(s) to simplify adding validation to structs. 


NOTE: this function is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterCustomTypeFunc](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L308)
                      [¶](#Validate.RegisterCustomTypeFunc)


                    [](#Validate)[](#CustomTypeFunc)
```
func (v *
                          Validate) RegisterCustomTypeFunc(fn 
                          CustomTypeFunc, types ...interface{})

```


RegisterCustomTypeFunc registers a CustomTypeFunc against a number of types 


NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterStructValidation](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L253)
                      [¶](#Validate.RegisterStructValidation)


                    [](#Validate)[](#StructLevelFunc)
```
func (v *
                          Validate) RegisterStructValidation(fn 
                          StructLevelFunc, types ...interface{})

```


RegisterStructValidation registers a StructLevelFunc against a number of types. 


NOTE: - this method is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterStructValidationCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L262)
                      [¶](#Validate.RegisterStructValidationCtx)


                    [](#Validate)[](#StructLevelFuncCtx)
```
func (v *
                          Validate) RegisterStructValidationCtx(fn 
                          StructLevelFuncCtx, types ...interface{})

```


RegisterStructValidationCtx registers a StructLevelFuncCtx against a number of types and allows passing of contextual validation information via context.Context. 


NOTE: - this method is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterStructValidationMapRules](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L281)
                      [¶](#Validate.RegisterStructValidationMapRules)


                      added in
                      v10.11.0


                    [](#Validate)[](/builtin#string)[](/builtin#string)
```
func (v *
                          Validate) RegisterStructValidationMapRules(rules map[
                          string]
                          string, types ...interface{})

```


RegisterStructValidationMapRules registers validate map rules. Be aware that map validation rules supersede those defined on a/the struct if present. 


NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterTagNameFunc](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L210)
                      [¶](#Validate.RegisterTagNameFunc)


                    [](#Validate)[](#TagNameFunc)
```
func (v *
                          Validate) RegisterTagNameFunc(fn 
                          TagNameFunc)

```


RegisterTagNameFunc registers a function to get alternate names for StructFields. 


eg. to use the names which have been specified for JSON representations of structs, rather than normal Go field names: 


```
validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    // skip if tag key says it should be ignored
    if name == "-" {
        return ""
    }
    return name
})
```


#### 
                    func (*Validate) [RegisterTranslation](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L321)
                      [¶](#Validate.RegisterTranslation)


                    [](#Validate)[](/builtin#string)[](/github.com/go-playground/universal-translator)[](/github.com/go-playground/universal-translator#Translator)[](#RegisterTranslationsFunc)[](#TranslationFunc)[](/builtin#error)
```
func (v *
                          Validate) RegisterTranslation(tag 
                          string, trans 
                          ut.
                          Translator, registerFn 
                          RegisterTranslationsFunc, translationFn 
                          TranslationFunc) (err 
                          error)

```


RegisterTranslation registers translations against the provided tag. 


#### 
                    func (*Validate) [RegisterValidation](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L220)
                      [¶](#Validate.RegisterValidation)


                    [](#Validate)[](/builtin#string)[](#Func)[](/builtin#bool)[](/builtin#error)
```
func (v *
                          Validate) RegisterValidation(tag 
                          string, fn 
                          Func, callValidationEvenIfNull ...
                          bool) 
                          error

```


RegisterValidation adds a validation with the given tag 


NOTES: - if the key already exists, the previous validation function will be replaced. - this method is not thread-safe it is intended that these all be registered prior to any validation 


#### 
                    func (*Validate) [RegisterValidationCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L226)
                      [¶](#Validate.RegisterValidationCtx)


                    [](#Validate)[](/builtin#string)[](#FuncCtx)[](/builtin#bool)[](/builtin#error)
```
func (v *
                          Validate) RegisterValidationCtx(tag 
                          string, fn 
                          FuncCtx, callValidationEvenIfNull ...
                          bool) 
                          error

```


RegisterValidationCtx does the same as RegisterValidation on accepts a FuncCtx validation allowing context.Context validation support. 


#### 
                    func (*Validate) [SetTagName](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L158)
                      [¶](#Validate.SetTagName)


                    [](#Validate)[](/builtin#string)
```
func (v *
                          Validate) SetTagName(name 
                          string)

```


SetTagName allows for changing of the default tag name of 'validate' 


#### 
                    func (*Validate) [Struct](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L345)
                      [¶](#Validate.Struct)


                    [](#Validate)[](/builtin#error)
```
func (v *
                          Validate) Struct(s interface{}) 
                          error

```


Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified. 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L354)
                      [¶](#Validate.StructCtx)


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#error)
```
func (v *
                          Validate) StructCtx(ctx 
                          context.
                          Context, s interface{}) (err 
                          error)

```


StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified and also allows passing of context.Context for contextual validation information. 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructExcept](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L522)
                      [¶](#Validate.StructExcept)


                    [](#Validate)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) StructExcept(s interface{}, fields ...
                          string) 
                          error

```


StructExcept validates all fields except the ones passed in. Fields may be provided in a namespaced fashion relative to the struct provided i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructExceptCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L533)
                      [¶](#Validate.StructExceptCtx)


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) StructExceptCtx(ctx 
                          context.
                          Context, s interface{}, fields ...
                          string) (err 
                          error)

```


StructExceptCtx validates all fields except the ones passed in and allows passing of contextual validation information via context.Context Fields may be provided in a namespaced fashion relative to the struct provided i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructFiltered](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L389)
                      [¶](#Validate.StructFiltered)


                    [](#Validate)[](#FilterFunc)[](/builtin#error)
```
func (v *
                          Validate) StructFiltered(s interface{}, fn 
                          FilterFunc) 
                          error

```


StructFiltered validates a structs exposed fields, that pass the FilterFunc check and automatically validates nested structs, unless otherwise specified. 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructFilteredCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L399)
                      [¶](#Validate.StructFilteredCtx)


                    [](#Validate)[](/context)[](/context#Context)[](#FilterFunc)[](/builtin#error)
```
func (v *
                          Validate) StructFilteredCtx(ctx 
                          context.
                          Context, s interface{}, fn 
                          FilterFunc) (err 
                          error)

```


StructFilteredCtx validates a structs exposed fields, that pass the FilterFunc check and automatically validates nested structs, unless otherwise specified and also allows passing of contextual validation information via context.Context 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructPartial](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L436)
                      [¶](#Validate.StructPartial)


                    [](#Validate)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) StructPartial(s interface{}, fields ...
                          string) 
                          error

```


StructPartial validates the fields passed in only, ignoring all others. Fields may be provided in a namespaced fashion relative to the struct provided eg. NestedStruct.Field or NestedArrayField[0].Struct.Name 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [StructPartialCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L447)
                      [¶](#Validate.StructPartialCtx)


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) StructPartialCtx(ctx 
                          context.
                          Context, s interface{}, fields ...
                          string) (err 
                          error)

```


StructPartialCtx validates the fields passed in only, ignoring all others and allows passing of contextual validation information via context.Context Fields may be provided in a namespaced fashion relative to the struct provided eg. NestedStruct.Field or NestedArrayField[0].Struct.Name 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. 


#### 
                    func (*Validate) [ValidateMap](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L194)
                      [¶](#Validate.ValidateMap)


                      added in
                      v10.6.0


                    [](#Validate)[](/builtin#string)[](/builtin#string)[](/builtin#string)
```
func (v *
                          Validate) ValidateMap(data map[
                          string]interface{}, rules map[
                          string]interface{}) map[
                          string]interface{}

```


ValidateMap validates map data from a map of tags 


#### 
                    func (Validate) [ValidateMapCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L164)
                      [¶](#Validate.ValidateMapCtx)


                      added in
                      v10.6.0


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#string)[](/builtin#string)[](/builtin#string)
```
func (v 
                          Validate) ValidateMapCtx(ctx 
                          context.
                          Context, data map[
                          string]interface{}, rules map[
                          string]interface{}) map[
                          string]interface{}

```


ValidateMapCtx validates a map using a map of validation rules and allows passing of contextual validation information via context.Context. 


#### 
                    func (*Validate) [Var](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L593)
                      [¶](#Validate.Var)


                    [](#Validate)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) Var(field interface{}, tag 
                          string) 
                          error

```


Var validates a single variable using tag style validation. eg. var i int validate.Var(i, "gt=1,lt=10") 


WARNING: a struct can be passed for validation eg. time.Time is a struct or if you have a custom type and have registered a custom type handler, so must allow it; however unforeseen validations will occur if trying to validate a struct that is meant to be passed to 'validate.Struct' 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. validate Array, Slice and maps fields which may contain more than one error 


#### 
                    func (*Validate) [VarCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L611)
                      [¶](#Validate.VarCtx)


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) VarCtx(ctx 
                          context.
                          Context, field interface{}, tag 
                          string) (err 
                          error)

```


VarCtx validates a single variable using tag style validation and allows passing of contextual validation information via context.Context. eg. var i int validate.Var(i, "gt=1,lt=10") 


WARNING: a struct can be passed for validation eg. time.Time is a struct or if you have a custom type and have registered a custom type handler, so must allow it; however unforeseen validations will occur if trying to validate a struct that is meant to be passed to 'validate.Struct' 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. validate Array, Slice and maps fields which may contain more than one error 


#### 
                    func (*Validate) [VarWithValue](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L646)
                      [¶](#Validate.VarWithValue)


                    [](#Validate)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) VarWithValue(field interface{}, other interface{}, tag 
                          string) 
                          error

```


VarWithValue validates a single variable, against another variable/field's value using tag style validation eg. s1 := "abcd" s2 := "abcd" validate.VarWithValue(s1, s2, "eqcsfield") // returns true 


WARNING: a struct can be passed for validation eg. time.Time is a struct or if you have a custom type and have registered a custom type handler, so must allow it; however unforeseen validations will occur if trying to validate a struct that is meant to be passed to 'validate.Struct' 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. validate Array, Slice and maps fields which may contain more than one error 


#### 
                    func (*Validate) [VarWithValueCtx](https://github.com/go-playground/validator/blob/v10.27.0/validator_instance.go#L665)
                      [¶](#Validate.VarWithValueCtx)


                    [](#Validate)[](/context)[](/context#Context)[](/builtin#string)[](/builtin#error)
```
func (v *
                          Validate) VarWithValueCtx(ctx 
                          context.
                          Context, field interface{}, other interface{}, tag 
                          string) (err 
                          error)

```


VarWithValueCtx validates a single variable, against another variable/field's value using tag style validation and allows passing of contextual validation information via context.Context. eg. s1 := "abcd" s2 := "abcd" validate.VarWithValue(s1, s2, "eqcsfield") // returns true 


WARNING: a struct can be passed for validation eg. time.Time is a struct or if you have a custom type and have registered a custom type handler, so must allow it; however unforeseen validations will occur if trying to validate a struct that is meant to be passed to 'validate.Struct' 


It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise. You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors. validate Array, Slice and maps fields which may contain more than one error 


#### 
                  type [ValidationErrors](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L36)
                    [¶](#ValidationErrors)


                  [](#FieldError)
```
type ValidationErrors []
                        FieldError

```


ValidationErrors is an array of FieldError's for use in custom error messages post validation. 


#### 
                    func (ValidationErrors) [Error](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L42)
                      [¶](#ValidationErrors.Error)


                    [](#ValidationErrors)[](/builtin#string)
```
func (ve 
                          ValidationErrors) Error() 
                          string

```


Error is intended for use in development + debugging and not intended to be a production error message. It allows ValidationErrors to subscribe to the Error interface. All information to create an error message specific to your application is contained within the FieldError found within the ValidationErrors array 


#### 
                    func (ValidationErrors) [Translate](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L54)
                      [¶](#ValidationErrors.Translate)


                    [](#ValidationErrors)[](/github.com/go-playground/universal-translator)[](/github.com/go-playground/universal-translator#Translator)[](#ValidationErrorsTranslations)
```
func (ve 
                          ValidationErrors) Translate(ut 
                          ut.
                          Translator) 
                          ValidationErrorsTranslations

```


Translate translates all of the ValidationErrors 


#### 
                  type [ValidationErrorsTranslations](https://github.com/go-playground/validator/blob/v10.27.0/errors.go#L17)
                    [¶](#ValidationErrorsTranslations)


                  [](/builtin#string)[](/builtin#string)
```
type ValidationErrorsTranslations map[
                        string]
                        string

```


ValidationErrorsTranslations is the translation return type
