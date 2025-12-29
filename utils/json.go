package utils

import (
	"encoding/json"

	"github.com/Verthandii/spring/it"
)

// JSONMarshal 对象转json字符串
func JSONMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// JSONMarshalByte 对象转字节
func JSONMarshalByte(v interface{}) []byte {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

// JSONDiff JSON对比器,返回前后两个json中不同的部分
func JSONDiff(s1, s2 []byte) ([]byte, []byte, error) {
	before, after, m1, m2 := it.H{}, it.H{}, it.H{}, it.H{}
	if len(s1) != 0 {
		if err := json.Unmarshal(s1, &m1); err != nil {
			return nil, nil, err
		}
	}
	if len(s2) != 0 {
		if err := json.Unmarshal(s2, &m2); err != nil {
			return nil, nil, err
		}
	}
	if err := getJSONDiff(m1, m2, before, after); err != nil {
		return nil, nil, err
	}
	if err := getJSONDiff(m2, m1, after, before); err != nil {
		return nil, nil, err
	}
	return JSONMarshalByte(before), JSONMarshalByte(after), nil
}

// getJSONDiff 获取两个对象结构差异
func getJSONDiff(m1, m2, before, after it.H) error {
	for k, m1Val := range m1 {
		if m2Val, ex := m2[k]; !ex {
			before[k] = m1Val
			after[k] = nil
		} else {
			switch m2Val.(type) {
			case int64, bool, float64, string:
				if m1Val != m2Val {
					before[k] = m1Val
					after[k] = m2Val
				}
			default:
				next1 := JSONMarshalByte(m1Val)
				next2 := JSONMarshalByte(m2Val)
				if string(next1) != string(next2) {
					// 判断下一层json差异
					nextDiff1, nextDiff2, err := JSONDiff(next1, next2)
					if err != nil {
						before[k] = m1Val
						after[k] = m2Val
					} else {
						mm1, mm2 := it.H{}, it.H{}
						if err := json.Unmarshal(nextDiff1, &mm1); err != nil {
							return err
						}
						if err := json.Unmarshal(nextDiff2, &mm2); err != nil {
							return err
						}
						before[k] = mm1
						after[k] = mm2
					}
				}
			}
		}
	}
	return nil
}
