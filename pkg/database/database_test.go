package database_test

import (
	"database/sql"
	"log"
	"testing"

	sideDb "github.com/jeefy/metrics-sidecar/pkg/database"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func TestMetricsUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sidecar Database Test")
}

func nodeMetrics() v1beta1.NodeMetricsList {
	tmp := v1beta1.NodeMetrics{}
	tmp.SetName("testing")
	tmp.Usage = v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse("10000"),
		v1.ResourceMemory: resource.MustParse("100000"),
	}

	nm := v1beta1.NodeMetricsList{
		Items: []v1beta1.NodeMetrics{
			tmp,
		},
	}

	return nm
}

func podMetrics() v1beta1.PodMetricsList {
	tmp2 := v1beta1.ContainerMetrics{}
	tmp2.Name = "container_test"
	tmp2.Usage = v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse("10000"),
		v1.ResourceMemory: resource.MustParse("100000"),
	}

	tmp := v1beta1.PodMetrics{}
	tmp.SetName("testing")
	tmp.Containers = []v1beta1.ContainerMetrics{
		tmp2,
	}

	nm := v1beta1.PodMetricsList{
		Items: []v1beta1.PodMetrics{
			tmp,
		},
	}

	return nm
}

var _ = Describe("Database functions", func() {
	Context("With an in-memory database", func() {
		It("should generate 'nodes' table to dump metrics in.", func() {
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				panic(err.Error())
			}

			defer db.Close()

			sideDb.CreateDatabase(db)

			_, err = db.Query("select * from nodes;")
			if err != nil {
				panic(err.Error())
			}
			Expect(err).To(BeNil())
		})

		It("should generate 'pods' table to dump metrics in.", func() {
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()

			sideDb.CreateDatabase(db)

			_, err = db.Query("select * from pods;")
			if err != nil {
				panic(err.Error())
			}
			Expect(err).To(BeNil())
		})

		It("should insert metrics into the database.", func() {
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()

			sideDb.CreateDatabase(db)

			nm := nodeMetrics()
			pm := podMetrics()

			sideDb.UpdateDatabase(db, &nm, &pm)

			rows, err := db.Query("select name, cpu, memory from nodes")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				var cpu string
				var memory string
				err = rows.Scan(&name, &cpu, &memory)
				if err != nil {
					log.Fatal(err)
				}
				Expect(err).To(BeNil())
				Expect(name).To(Equal("testing"))
				Expect(cpu).To(Equal("10k"))
				Expect(memory).To(Equal("100k"))
			}

			rows, err = db.Query("select name, container, cpu, memory from pods")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				var container string
				var cpu string
				var memory string
				err = rows.Scan(&name, &container, &cpu, &memory)
				if err != nil {
					log.Fatal(err)
				}
				Expect(err).To(BeNil())
				Expect(name).To(Equal("testing"))
				Expect(container).To(Equal("container_test"))
				Expect(cpu).To(Equal("10k"))
				Expect(memory).To(Equal("100k"))
			}
		})
		It("should insert metrics into the database.", func() {
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()

			sideDb.CreateDatabase(db)

			nm := nodeMetrics()
			pm := podMetrics()

			sideDb.UpdateDatabase(db, &nm, &pm)

			sqlStmt := "insert into nodes(name,cpu,memory,storage,time) values('lame','20k','300k','0',datetime('now','-20 minutes'));"
			_, err = db.Exec(sqlStmt)
			if err != nil {
				panic(err.Error())
			}

			timeWindow := 5
			sideDb.CullDatabase(db, &timeWindow)

			rows, err := db.Query("select name, cpu, memory from nodes")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				var cpu string
				var memory string
				err = rows.Scan(&name, &cpu, &memory)
				if err != nil {
					log.Fatal(err)
				}
				Expect(err).To(BeNil())
				Expect(name).To(Equal("testing"))
				Expect(cpu).To(Equal("10k"))
				Expect(memory).To(Equal("100k"))
			}

		})
	})
})
