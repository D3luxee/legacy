package legacy // import "github.com/solnx/legacy"

import (
	"fmt"

	"github.com/raintank/schema"
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
		metric.Name = m.Path
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
			tag := string(k) + "=" + v
			metric.Tags = append(metric.Tags, tag)
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
