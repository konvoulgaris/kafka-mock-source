# kafka-mock-source

A simple tool to mock a source that generates and streams data to Kafka. It supports flexible generation of fake data based on faker and other configurations using a YAML file. This is mainly intended to be used for development purposes.

## How to Run

### 1. Compile the Program

```bash
go build
```

### 2. Copy it to the System Path
```bash
sudo cp kafka-mock-source /usr/local/bin
```

### 3. Run it
```bash
kafka-mock-source <path-to-config.yaml>
```

## Configuration (YAML)

### Example Configuration File for Stock Price Scenario
This configuration file, along with the `symbols.json` provided in the repository, generates mock stock data with:
+ `previousPrice`: Float value representing the previous price of the stock.
+ `currentPrice`: Float value representing the current price of the stock.
+ `marketConfidence`: Factor value representing the market confidence.
+ `latestHeadline`: Random sentence representing the latest headline regarding the stock at that time.
+ `timestamp`: The time that the stock buy/sale was recorded.
+ `ticker`: From the `correlation` settings, the ticker symbol of the stock. Sourced from the `symbols.json` file.


```yaml
kafka: localhost:9092
topic: test
interval: 1000
samples: 5
format: csvheadless
correlation:
  label: ticker
  source: ./symbols.json
data:
  - previousPrice=randomfloat
  - currentPrice=randomfloat
  - marketConfidence=randomfactor
  - latestHeadline=sentence
  - timestamp=timestamp
```

### Structure

+ `kafka`: Sets the host and port of the Kafka server.
  - **Options**: `<host>:<port>`

+ `topic`: Defines the Kafka topic where the data will be published.
  - **Options**: String

+ `interval`: Specifies the time interval (in milliseconds) between publishing each payload.
  - **Options**: Positive Integers

+ `samples`: Sets the number of samples to be generated and published per payload.
  - **Options**: Positive Integers

+ `format`: Specifies the format of the data to be published (json, csv, or csvheadless).
  - **Options**: json, csv, csvheadless

+ `correlation`: Defines the correlation settings for generating correlated data.
  - `label`: Specifies the data field label that will be appened to the generated data.
    - **Options**: String
  - `amount`: Defines the size of the pool of correlation IDs available for selection. If set to 5, each sample can correlate with any of those 5 IDs.
    - **Options**: Positive Integers
  - `source`: Defines a source file from where we get the correlation IDs from, instead of generating them randomly.
    - **Options**: A path to valid JSON file of an array of strings.

+ `data`: Defines the types of data to be generated for each entry.
  - **Options**: List of key-value pairs

**NOTE:** `correlation.amount` and `correlation.source` cannot be set at the same time.

### Supported Data Types

| Fake Type           | Description                                                  |
|---------------------|--------------------------------------------------------------|
| amountwithcurrency  | Generates an amount with currency symbol.                   |
| ccnumber            | Generates a credit card number.                              |
| cctype              | Generates a credit card type.                                |
| century             | Generates a century.                                         |
| chinesefirstname    | Generates a Chinese first name.                              |
| chineselastname     | Generates a Chinese last name.                               |
| chinesename         | Generates a Chinese name.                                    |
| currency            | Generates a currency symbol.                                 |
| date                | Generates a date.                                            |
| dayofmonth          | Generates a day of the month.                                |
| dayofweek           | Generates a day of the week.                                 |
| domainname          | Generates a domain name.                                     |
| e164phonenumber     | Generates an E164 phone number.                              |
| email               | Generates an email address.                                  |
| firstname           | Generates a first name.                                      |
| firstnamefemale     | Generates a female first name.                               |
| firstnamemale       | Generates a male first name.                                 |
| gender              | Generates a gender.                                          |
| ipv4                | Generates an IPv4 address.                                   |
| ipv6                | Generates an IPv6 address.                                   |
| jwt                 | Generates a JSON Web Token (JWT).                            |
| lastname            | Generates a last name.                                       |
| latitude            | Generates a latitude.                                        |
| longitude           | Generates a longitude.                                       |
| macaddress          | Generates a MAC address.                                     |
| monthname           | Generates a month name.                                      |
| name                | Generates a name.                                            |
| paragraph           | Generates a paragraph of text.                               |
| password            | Generates a password.                                        |
| phonenumber         | Generates a phone number.                                    |
| randomint           | Generates a random integer in the 0 to 100 range.            |
| randomfloat         | Generates a random float in the 0 to 100 range.              |
| randomfactor        | Generates a random factor.                                   |
| randomunixtime      | Generates a random Unix time.                                |
| sentence            | Generates a sentence.                                        |
| timestring          | Generates a time string.                                     |
| timeperiod          | Generates a time period.                                     |
| timestamp           | Generates a timestamp.                                       |
| timezone            | Generates a timezone.                                        |
| titlefemale         | Generates a female title.                                    |
| titlemale           | Generates a male title.                                      |
| tollfreephonenumber | Generates a toll-free phone number.                          |
| url                 | Generates a URL.                                             |
| uuiddigit           | Generates a UUID digit.                                      |
| uuidhyphenated      | Generates a hyphenated UUID.                                 |
| unixtime            | Generates a Unix time.                                       |
| username            | Generates a username.                                        |
| word                | Generates a word.                                            |
| yearstring          | Generates a year string.                                     |

## License

The code in this repository is licensed under the [Apache Licence Version 2.0](LICENSE) by [Konstantinos Voulgaris](https://github.com/konvoulgaris).
