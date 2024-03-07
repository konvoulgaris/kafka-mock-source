package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"gopkg.in/yaml.v3"
)

var fakerMap = map[string]func(...options.OptionFunc) string{
	"amountwithcurrency":  faker.AmountWithCurrency,
	"ccnumber":            faker.CCNumber,
	"cctype":              faker.CCType,
	"century":             faker.Century,
	"chinesefirstname":    faker.ChineseFirstName,
	"chineselastname":     faker.ChineseLastName,
	"chinesename":         faker.ChineseName,
	"currency":            faker.Currency,
	"date":                faker.Date,
	"dayofmonth":          faker.DayOfMonth,
	"dayofweek":           faker.DayOfWeek,
	"domainname":          faker.DomainName,
	"e164phonenumber":     faker.E164PhoneNumber,
	"email":               faker.Email,
	"firstname":           faker.FirstName,
	"firstnamefemale":     faker.FirstNameFemale,
	"firstnamemale":       faker.FirstNameMale,
	"gender":              faker.Gender,
	"ipv4":                faker.IPv4,
	"ipv6":                faker.IPv6,
	"jwt":                 faker.Jwt,
	"lastname":            faker.LastName,
	"latitude":            func(...options.OptionFunc) string { return strconv.FormatFloat(faker.Latitude(), 'f', -1, 64) },
	"longitude":           func(...options.OptionFunc) string { return strconv.FormatFloat(faker.Longitude(), 'f', -1, 64) },
	"macaddress":          faker.MacAddress,
	"monthname":           faker.MonthName,
	"name":                faker.Name,
	"paragraph":           faker.Paragraph,
	"password":            faker.Password,
	"phonenumber":         faker.Phonenumber,
	"randomint":           func(...options.OptionFunc) string { return strconv.Itoa(rand.Intn(101)) },
	"randomunixtime":      func(...options.OptionFunc) string { return strconv.FormatInt(faker.RandomUnixTime(), 10) },
	"sentence":            faker.Sentence,
	"timestring":          faker.TimeString,
	"timeperiod":          faker.Timeperiod,
	"timestamp":           faker.Timestamp,
	"timezone":            faker.Timezone,
	"titlefemale":         faker.TitleFemale,
	"titlemale":           faker.TitleMale,
	"tollfreephonenumber": faker.TollFreePhoneNumber,
	"url":                 faker.URL,
	"uuiddigit":           faker.UUIDDigit,
	"uuidhyphenated":      faker.UUIDHyphenated,
	"unixtime":            func(...options.OptionFunc) string { return strconv.FormatInt(faker.RandomUnixTime(), 10) },
	"username":            faker.Username,
	"word":                faker.Word,
	"yearstring":          faker.YearString,
}

type Config struct {
	Data []string `yaml:"data"`
}

type DataField struct {
	Label string
	Value string
}

func getDataField(line string) DataField {
	data := strings.Split(line, "=")

	if len(data) != 2 {
		log.Println("Skipping invalid data field:", line)
		return DataField{}
	}

	return DataField{
		Label: data[0],
		Value: data[1],
	}
}

func main() {
	configFile, err := os.ReadFile("config.yaml")

	if err != nil {
		println("Error opening config file!")
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		println("Error parsing config file!")
		panic(err)
	}

	var dataFields []DataField

	for _, line := range config.Data {
		dataFields = append(dataFields, getDataField(line))
	}

	for _, dataField := range dataFields {
		key := dataField.Value

		if _, ok := fakerMap[key]; !ok {
			log.Println("Skipping invalid data field value:", key)
			continue
		}

		value := fakerMap[key]()

		println(dataField.Label, "=", value)
	}
}
