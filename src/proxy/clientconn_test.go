package proxy

import (
	"testing"
)

func newTestClient() *ProxyConn {
	c := NewProxyConn()

	if err := c.Connect("127.0.0.1:3306", "qing", "admin", "mixer"); err != nil {
		return nil
	}

	return c
}

func TestClientConn_Handshake(t *testing.T) {
	newTestServer()

	c := newTestClient()
	if c == nil {
		t.Fatal("connect failed")
	}

	c.Close()
}

func TestClientConn_CreateTable(t *testing.T) {
	s := `CREATE TABLE IF NOT EXISTS mixer_test_clientconn (
          id BIGINT(64) UNSIGNED  NOT NULL,
          str VARCHAR(256),
          f DOUBLE,
          e enum("test1", "test2"),
          PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	c := newTestProxyConn()

	if _, err := c.Exec(s); err != nil {
		t.Fatal(err)
	}

}

func TestClientConn_Delete(t *testing.T) {
	s := `delete from mixer_test_clientconn`

	c := newTestClient()
	defer c.Close()

	_, err := c.Exec(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientConn_Insert(t *testing.T) {
	s := `insert into mixer_test_clientconn (id, str, f, e) values (1, "abc", 3.14, "test1")`

	c := newTestClient()
	defer c.Close()

	pkg, err := c.Exec(s)
	if err != nil {
		t.Fatal(err)
	}

	if pkg.AffectedRows != 1 {
		t.Fatal(pkg.AffectedRows)
	}
}

func TestClientConn_Select(t *testing.T) {
	s := `select str, f, e from mixer_test_clientconn where id = 1`

	c := newTestClient()
	defer c.Close()

	result, err := c.Query(s)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.ColumnDefs) != 3 {
		t.Fatal(len(result.ColumnDefs))
	}

	if len(result.Rows) != 1 {
		t.Fatal(len(result.Rows))
	}
}

func TestClientConn_DeleteTable(t *testing.T) {
	c := newTestProxyConn()

	if _, err := c.Exec("drop table mixer_test_clientconn"); err != nil {
		t.Fatal(err)
	}
}
