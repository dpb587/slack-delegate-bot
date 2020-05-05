package configutil_test

import (
	. "github.com/dpb587/slack-delegate-bot/pkg/configutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemarshalYaml", func() {
	It("remarshals", func() {
		dst := struct {
			First  string   `yaml:"one"`
			Second []string `yaml:"two"`
		}{}

		err := RemarshalYAML(
			map[string]interface{}{"one": "primary", "two": []string{"secondary", "tertiary"}},
			&dst,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(dst.First).To(Equal("primary"))
		Expect(dst.Second).To(Equal([]string{"secondary", "tertiary"}))
	})

	It("unmarshals strictly", func() {
		dst := struct {
			First string `yaml:"one"`
		}{}

		err := RemarshalYAML(map[string]interface{}{"two": "missing"}, &dst)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unmarshalling"))
	})
})
