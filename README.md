Masq - Time-sensitive key/value storage
=======================================

# API Usage


## Generating Passwords (/passwords)

Masquerade exposes an endpoint that enables simple password generation. To utilize this, make a _GET_ request to:

```
https://www.example.org/passwords
```


Which returns the following _JSON_

```
{
  "password": "VomitousCrimple1552",
}
```


### Response Properties

* `password` - [String] The pre-generated password, using _<word[0]> + <word[1]> + padding_

### Warning
The password generator uses a dictionary to generate the passwords. It is not intended for long-term passwords but, limited-use such as default passwords when provisioning accounts



## Storing Information (/api/hide)

To _hide_ information in masquerade, an HTTP _POST_ request must be made to:

```
https://www.example.org/hide?data=<value to store>&duration=<duration to store for>
```

### Required Properties

* `data` - [String] The information to store in Masquerade.  return


### Optional Properties

* `duration` - [String] The duration to store the information for
    * Must be one of:
        * `5m` (5 minutes)
        * `15m` (15 minutes)
        * `30m` (30 minutes)
        * `1h` (1 hour)
        * `24h` (Default; 24 Hours)
        * `48h` (48 Hours)
        * `72h` (72 Hours)
        * `1w` (1 Week)
    * _If no duration is given, masquerade defaults to storing the information for 24 hours_  return


### Output

_The output is in JSON format_

* `url` - [String] The browser/user-friendly URL which can be used to retrieve the stored information
* `key` - [String] The hash key which the data is stored under
* `duration` - [String] The duration that the information is stored for  return


## Retrieving Information (/show)

To `retrieve` information from masquerade, an HTTP `get` request must be made to:

```
https://www.example.org/api/show?key=<hash key>
```


### Required URL Parameters

* `key` - [String] This is the hash key the value is stored under


### Output

_The output is in JSON format_

* `value` - [String] The stored information


# Examples


## Storing the value `test` for `72` hours:

```bash
curl --insecure -s -X POST "https://www.example.org/hide?data=test&duration=72h" -H 'Content-type:application/json'
```

will return:

```
{
  "duration": "72h",
  "key": "1eAf5b",
  "url": "https://www.example.org/s/1eAf5b"
}
```


## Retrieving a stored value:

```
curl --insecure -s "https://www.example.org/show?key=1eAf5b"
```

will return:

```
{
  "value": "test"
}
```

