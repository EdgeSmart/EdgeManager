package dao

import "fmt"

func GetClusterData(clusterKey string) map[string]string {
	db, _ := GetDB("edge")

	rows, err := db.Query("SELECT `cid`, `ckey`, `uid`, `name`, `token`, `status` FROM `cluster` WHERE ckey = ?", clusterKey)
	if err != nil {
		fmt.Println(err)
		return map[string]string{}
	}
	data := map[string]string{}
	for rows.Next() {
		var (
			cid    string
			ckey   string
			uid    string
			name   string
			token  string
			status string
		)
		err = rows.Scan(&cid, &ckey, &uid, &name, &token, &status)
		data = map[string]string{
			"cid":    cid,
			"ckey":   ckey,
			"uid":    uid,
			"name":   name,
			"token":  token,
			"status": status,
		}
	}
	return data
}
