package legacy // import "github.com/solnx/legacy"

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/d3luxee/schema"
)

func (m *MetricSplit) ConvertToMetrics20() (schema.MetricData, error) {
	if m.IsMetrics20() {
		if m.Type == "string" {
			return schema.MetricData{}, fmt.Errorf("string metrics can't be converted")
		}
		metric := schema.MetricData{}
		//Use public OrgID
		metric.OrgId = 1
		//Use the path as name
		metric.Name = strings.Replace(m.Path, "/", ".", -1)
		//All values are float64 in metrics20
		switch m.Type {
		case `integer`:
			fallthrough
		case `long`:
			metric.Value = float64(m.Val.IntVal)
		case `real`:
			metric.Value = m.Val.FlpVal
		}
		//always asume the mtype gauge because the source does not provide the necessary information
		metric.Mtype = "gauge"
		metric.Interval = 10
		metric.Time = m.TS.Unix()
		//We do not know the unit
		metric.Unit = ""
		//Add tags
		for k, v := range m.Tags {
			// The old tags where an array to keep them we will convert them to key:value pairs
			// using the index of the array as key
			if !isUUID(v) {
				tag := "what=" + v
				metric.Tags = append(metric.Tags, tag)
			} else {
				tag := string(k) + "=" + v
				metric.Tags = append(metric.Tags, tag)
			}
		}
		for k, v := range m.Labels {
			tag := k + "=" + v
			metric.Tags = append(metric.Tags, tag)
		}
		metric.SetId()
		return metric, metric.Validate()
	}
	return schema.MetricData{}, fmt.Errorf("this metric is not intended to be a metrics20 metric")
}

// isUUID validates if a string is one very narrow formatting of a
// UUID. Other valid formats with braces etc are not accepted.
func isUUID(s string) bool {
	const reUUID string = `^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[1-5][[:xdigit:]]{3}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`
	const reUNIL string = `^0{8}-0{4}-0{4}-0{4}-0{12}$`
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", reUUID, reUNIL))

	return re.MatchString(s)
}
