package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

var correlationIds []string

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
	"randomfloat":         func(...options.OptionFunc) string { return strconv.FormatFloat(rand.Float64()*100, 'f', -1, 64) },
	"randomfactor":        func(...options.OptionFunc) string { return strconv.FormatFloat(rand.Float64(), 'f', -1, 64) },
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

type ConfigCorrelation struct {
	Amount int    `yaml:"amount"`
	Label  string `yaml:"label"`
}

type Config struct {
	Kafka       string            `yaml:"kafka"`
	Topic       string            `yaml:"topic"`
	Interval    int               `yaml:"interval"`
	Samples     int               `yaml:"samples"`
	Format      string            `yaml:"format"`
	Correlation ConfigCorrelation `yaml:"correlation"`
	Data        []string          `yaml:"data"`
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

func generateSample(dataFields []DataField) map[string]interface{} {
	data := make(map[string]interface{})

	for _, dataField := range dataFields {
		fakeKey := dataField.Value

		if fakeKey == "correlate" {
			data[dataField.Label] = correlationIds[rand.Intn(len(correlationIds))]
			continue
		}

		data[dataField.Label] = fakerMap[fakeKey]()
	}

	return data
}

func generateCSVData(dataFields []DataField, data []map[string]interface{}, headless bool) []byte {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if !headless {
		var headers []string

		for _, dataField := range dataFields {
			headers = append(headers, dataField.Label)
		}

		err := writer.Write(headers)

		if err != nil {
			log.Println("Error writing CSV headers:", err)
		}
	}

	for _, row := range data {
		sample := make([]string, len(dataFields))

		for i, dataField := range dataFields {
			sample[i] = row[dataField.Label].(string)
		}

		err := writer.Write(sample)

		if err != nil {
			log.Println("Error writing CSV row:", err)
		}
	}

	writer.Flush()
	return buffer.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run main.go <path-to-config.yaml>")
	}

	configFile, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatalln("Error reading config file:", err)
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		log.Fatalln("Error parsing config file:", err)
	}

	if len(config.Kafka) <= 0 {
		log.Fatalln("Invalid kafka value:", config.Kafka, "\nValid values are non-empty strings")
	}

	if len(config.Topic) <= 0 {
		log.Fatalln("Invalid topic value:", config.Topic, "\nValid values are non-empty strings")
	}

	if config.Interval <= 0 {
		log.Fatalln("Invalid interval value:", config.Interval, "\nValid values are positive integers")
	}

	if config.Samples <= 0 {
		log.Fatalln("Invalid samples value:", config.Samples, "\nValid values are positive integers")
	}

	if config.Format == "" || (config.Format != "json" && config.Format != "csv" && config.Format != "csvheadless") {
		log.Fatalln("Invalid format value:", config.Format, "\nValid values are: 'json', 'csv', 'csvheadless'")
	}

	var dataFields []DataField

	for _, line := range config.Data {
		temp := getDataField(line)

		if _, ok := fakerMap[temp.Value]; !ok {
			log.Fatalln("Invalid data field value:", temp.Value, "\nSee README.md for a list of valid values")
		}

		if config.Correlation.Label == temp.Label {
			log.Fatalln("Invalid data field label:", temp.Label, "\nThat label is a reserved name used for data correlation based on your config.\nPlease use a different label.")
		}

		dataFields = append(dataFields, temp)
	}

	if len(dataFields) <= 0 {
		log.Fatalln("No valid data fields found in config file")
	}

	if config.Correlation.Amount <= 0 && config.Correlation.Label != "" {
		log.Fatalln("Invalid correlation amount value:", config.Correlation.Amount, "\nValid values are positive integers")
	} else if config.Correlation.Amount > 0 && config.Correlation.Label == "" {
		log.Fatalln("Invalid correlation label value:", config.Correlation.Label, "\nValid values are non-empty strings")
	} else {
		dataFields = append(dataFields, DataField{
			Label: config.Correlation.Label,
			Value: "correlate",
		})

		for i := 0; i < config.Correlation.Amount; i++ {
			correlationIds = append(correlationIds, faker.UUIDDigit())
		}
	}

	client, err := kafka.DialLeader(context.Background(), "tcp", config.Kafka, config.Topic, 0)

	if err != nil {
		log.Fatalln("Error connecting to Kafka:", err)
	}

	defer client.Close()

	for {
		var payload []map[string]interface{}

		for i := 0; i < config.Samples; i++ {
			data := generateSample(dataFields)
			payload = append(payload, data)
		}

		var message []byte

		if config.Format == "json" {
			jsonData, err := json.Marshal(payload)

			if err != nil {
				log.Println("Error marshalling messages to JSON:", err)
			}

			message = jsonData
		} else if config.Format == "csv" {
			message = generateCSVData(dataFields, payload, false)
		} else {
			message = generateCSVData(dataFields, payload, true)
		}

		n, err := client.Write(message)

		if err != nil {
			log.Println("Error writing messages to Kafka:", err)
		}

		log.Println("Wrote", n, "bytes to Kafka")
		time.Sleep(time.Duration(config.Interval) * time.Millisecond)
	}
}
