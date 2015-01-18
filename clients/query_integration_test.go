package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/tidepool-org/go-common/clients/mongo"
	"labix.org/v2/mgo"

	"../model"
)

const (
	QUERY_WHERE_AND = "METAQUERY WHERE userid IS 1234 QUERY TYPE IN settings WHERE time >= 2014-10-23T07:00:00.000Z AND time < 2014-10-23T08:00:00.000Z SORT BY time AS Timestamp REVERSED"
	QUERY_WHERE     = "METAQUERY WHERE userid IS 1234 QUERY TYPE IN basal WHERE time <= 2014-10-23T08:00:00.000Z SORT BY time AS Timestamp REVERSED"
	QUERY_WHERE_IN  = "METAQUERY WHERE userid IS 1234 QUERY TYPE IN basal, settings WHERE uploadId NOT IN test-data3, test-data2 SORT BY time AS Timestamp REVERSED"
)

func setupForTest() *MongoStoreClient {
	//we are setting to something other than the default so we can isolate the test data
	testingConfig := &mongo.Config{ConnectionString: "mongodb://localhost/streams_test"}

	mc := NewMongoStoreClient(testingConfig)

	/*
	 * INIT THE TEST - we use a clean copy of the collection before we start
	 */

	mc.session.DB("").C(DEVICE_DATA_COLLECTION).DropCollection()

	if err := mc.session.DB("").C(DEVICE_DATA_COLLECTION).Create(&mgo.CollectionInfo{}); err != nil {
		log.Panicf("We couldn't create the device data collection for these tests ", err)
	}

	/*
	 * Load test data
	 */
	if testData, err := ioutil.ReadFile("./test_data.json"); err == nil {

		var toLoad []interface{}

		if err := json.Unmarshal(testData, &toLoad); err != nil {
			log.Panicf("We could not load the test data ", err)
		}

		for i := range toLoad {
			//insert each test data item
			if insertErr := mc.session.DB("").C(DEVICE_DATA_COLLECTION).Insert(toLoad[i]); insertErr != nil {
				log.Panicf("We could not load the test data ", insertErr)
			}
		}
	}

	return mc

}

func Test_Full_WithWhereAnd(t *testing.T) {

	//lets test it all

	errs, qd := model.BuildQuery(QUERY_WHERE_AND)

	if len(errs) != 0 {
		t.Fatalf("there should be no errors but got %v", errs)
	}

	mc := setupForTest()

	if results := mc.ExecuteQuery(qd); results == nil {
		t.Fatalf("no results were found for the query [%v]", qd)
	} else {

		type found map[string]interface{}

		records := []found{}
		json.Unmarshal(results, &records)

		if len(records) != 1 {
			t.Fatalf("there should only be one result but got [%d]", len(records))
		}

		if records[0]["type"] != "settings" {
			t.Fatalf("should have been a settings record %v", records[0])
		}

	}

}

func Test_Full_WithWhereOnly(t *testing.T) {

	//lets test it all

	errs, qd := model.BuildQuery(QUERY_WHERE)

	if len(errs) != 0 {
		t.Fatalf("there should be no errors but got %v", errs)
	}

	mc := setupForTest()

	if results := mc.ExecuteQuery(qd); results == nil {
		t.Fatalf("no results were found for the query [%v]", qd)
	} else {
		type found map[string]interface{}

		records := []found{}
		json.Unmarshal(results, &records)

		if len(records) != 2 {
			t.Fatalf("there should be 2 results but got [%d]", len(records))
		}

		if records[0]["type"] != "basal" {
			t.Fatalf("should have been a basal record %v", records[0])
		}

		if records[1]["type"] != "basal" {
			t.Fatalf("should have been a basal record %v", records[1])
		}
	}

}

func Test_Full_WithWhereIn(t *testing.T) {

	errs, qd := model.BuildQuery(QUERY_WHERE_IN)

	if len(errs) != 0 {
		t.Fatalf("there should be no errors but got %v", errs)
	}

	mc := setupForTest()

	if results := mc.ExecuteQuery(qd); results == nil {
		t.Fatalf("no results were found for the query [%v]", qd)
	} else {
		type found map[string]interface{}

		records := []found{}
		json.Unmarshal(results, &records)

		if len(records) != 5 {
			t.Fatalf("there should be 5 results but got [%d]", len(records))
		}

	}

}
