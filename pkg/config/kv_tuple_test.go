package config_test

import (
	. "github.com/dpb587/slack-delegate-bot/pkg/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("KvTuple", func() {
	DescribeTable(
		"requires a map input",
		func(input interface{}) {
			_, _, err := KeyValueTuple(input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected map[string]interface{}"))
		},
		Entry("string", "fake-data"),
		Entry("expected map[int]interface", map[int]interface{}{}),
		Entry("int", 123),
	)

	It("errors if non-string keys are present", func() {
		_, _, err := KeyValueTuple(map[interface{}]interface{}{
			"key1": nil,
			123:    nil,
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected string key"))
	})

	It("errors if multiple keys are present", func() {
		_, _, err := KeyValueTuple(map[interface{}]interface{}{
			"key1": nil,
			"key2": nil,
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected exactly one key-value tuple"))
	})

	It("errors if no keys are present", func() {
		_, _, err := KeyValueTuple(map[interface{}]interface{}{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected exactly one key-value tuple"))
	})

	It("works", func() {
		key, value, err := KeyValueTuple(map[interface{}]interface{}{
			"fake-key1": []string{"fake-value1"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(key).To(Equal("fake-key1"))
		Expect(value).To(Equal([]string{"fake-value1"}))
	})
})
