package translation_test

import (
	"bytes"
	"errors"
	"fmt"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	settranslationfakes "github.com/jwfriese/workouttrackerapi/lifts/sets/translation/translationfakes"
	lifttranslation "github.com/jwfriese/workouttrackerapi/lifts/translation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LiftsCreateRequestTranslator", func() {
	var (
		subject           lifttranslation.LiftsCreateRequestTranslator
		fakeSetTranslator *settranslationfakes.FakeSetsCreateRequestTranslator
	)

	BeforeEach(func() {
		fakeSetTranslator = new(settranslationfakes.FakeSetsCreateRequestTranslator)
		subject = lifttranslation.NewLiftsCreateRequestTranslator(fakeSetTranslator)
	})

	Describe("Translating JSON into a lift", func() {
		var (
			result *liftdatamodel.Lift
			err    error
		)

		Context("When given valid JSON", func() {
			var (
				setOne *setdatamodel.Set
				setTwo *setdatamodel.Set
			)

			BeforeEach(func() {
				setOne = &setdatamodel.Set{}
				setTwo = &setdatamodel.Set{}
				fakeSetTranslator.TranslateStub = func(requestJSON []byte) (*setdatamodel.Set, error) {
					if bytes.Equal([]byte(`{"name":"setOne"}`), requestJSON) {
						return setOne, nil
					} else if bytes.Equal([]byte(`{"name":"setTwo"}`), requestJSON) {
						return setTwo, nil
					}

					errString := fmt.Sprintf("Invalid requestJSON (%q). Expected (%q) or (%q)", requestJSON, []byte(`{"name":"setOne}"`), []byte(`{"name":"setTwo"}`))
					return nil, errors.New(errString)
				}

				validJSON := []byte(`{"name":"turtle shoulder press","dataTemplate":"weight/time_in_seconds", "sets":[{"name":"setOne"},{"name":"setTwo"}]}`)

				result, err = subject.Translate(validJSON)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("translates the JSON into a lift", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("turtle shoulder press"))
				Expect(result.DataTemplate).To(Equal("weight/time_in_seconds"))
				Expect(len(result.Sets)).To(Equal(2))
				Expect(result.Sets[0]).To(BeIdenticalTo(setOne))
				Expect(result.Sets[1]).To(BeIdenticalTo(setTwo))
			})
		})

		Context("When given JSON with an error", func() {
			Context("Missing 'name'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"dataTemplate":"weight/time_in_seconds", "sets":[{"name":"setOne"},{"name":"setTwo"}]}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns a descriptive error", func() {
					Expect(result).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Missing required \"name\" field in lift request JSON"))
				})
			})

			Context("Missing 'dataTemplate'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"name":"turtle shoulder press", "sets":[{"name":"setOne"},{"name":"setTwo"}]}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns a descriptive error", func() {
					Expect(result).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Missing required \"dataTemplate\" field in lift request JSON"))
				})
			})

			Context("Missing 'sets'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"name":"turtle shoulder press", "dataTemplate":"weight/time_in_seconds"}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns a descriptive error", func() {
					Expect(result).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Missing required \"sets\" field in lift request JSON"))
				})
			})
		})
	})
})
