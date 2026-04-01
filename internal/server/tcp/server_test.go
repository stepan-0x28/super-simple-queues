package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"super-simple-queues/internal/queue"
	"testing"
	"time"
)

var (
	queueManager *queue.Manager
	tcpPort      int
)

func TestMain(m *testing.M) {
	const (
		queueChunkSize    = 1024
		tcpConnBufferSize = 256
	)

	queueManager = queue.NewManager(queueChunkSize)

	tcpServer := NewServer(queueManager, tcpConnBufferSize)

	go func() { _ = tcpServer.Run(0) }()

	time.Sleep(time.Second)

	tcpPort = tcpServer.listener.Addr().(*net.TCPAddr).Port

	code := m.Run()

	_ = tcpServer.listener.Close()

	os.Exit(code)
}

func TestServer_sendingQueueItems(t *testing.T) {
	const (
		queueKey                = "test"
		expectedQueueItemsCount = 64
	)

	createQueue(t, queueKey)

	sendQueueItems(t, queueKey, expectedQueueItemsCount)

	checkExpectedQueueItemsCount(t, queueKey, expectedQueueItemsCount)
}

func TestServer_receivingQueueItems(t *testing.T) {
	const (
		queueKey                = "test"
		queueItemsSentCount     = 64
		queueItemsReceivedCount = 16
		expectedQueueItemsCount = queueItemsSentCount - queueItemsReceivedCount
	)

	createQueue(t, queueKey)

	sendQueueItems(t, queueKey, queueItemsSentCount)

	conn := connectTCPServer(t)

	sendMessage(t, conn, generateInitMessage(0, queueKey))

	for i := 0; i < queueItemsReceivedCount; i++ {
		receiveQueueItem(t, conn)
	}

	_ = conn.Close()

	checkExpectedQueueItemsCount(t, queueKey, expectedQueueItemsCount)
}

func TestServer_dataAccuracy(t *testing.T) {
	const (
		queueKey              = "test"
		expectedQueueItemData = "queue_item_data"
	)

	createQueue(t, queueKey)

	conn := connectTCPServer(t)

	sendMessage(t, conn, generateInitMessage(1, queueKey))
	sendMessage(t, conn, generatePayloadMessage(expectedQueueItemData))

	_ = conn.Close()

	conn = connectTCPServer(t)

	sendMessage(t, conn, generateInitMessage(0, queueKey))

	queueItem := string(receiveQueueItem(t, conn))

	if queueItem != expectedQueueItemData {
		t.Fatalf("expected data for queue item %v, received %v", expectedQueueItemData, queueItem)
	}
}

func createQueue(t *testing.T, queueKey string) {
	queueManager.Create(queueKey)

	t.Cleanup(func() { queueManager.Delete(queueKey) })
}

func sendQueueItems(t *testing.T, queueKey string, itemsCount int) {
	t.Helper()

	conn := connectTCPServer(t)

	sendMessage(t, conn, generateInitMessage(1, queueKey))

	const queueItemData = "queue_item_data"

	for i := 0; i < itemsCount; i++ {
		sendMessage(t, conn, generatePayloadMessage(queueItemData))
	}

	_ = conn.Close()
}

func connectTCPServer(t *testing.T) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", tcpPort))

	if err != nil {
		t.Fatalf("error connecting to the TCP server, %v", err)
	}

	return conn
}

func generateInitMessage(operatingMode byte, queueKey string) []byte {
	initMessage := []byte{1, operatingMode, byte(len(queueKey))}

	initMessage = append(initMessage, []byte(queueKey)...)

	return initMessage
}

func generatePayloadMessage(queueItemData string) []byte {
	payloadMessage := []byte{2}

	payloadMessage = binary.BigEndian.AppendUint32(payloadMessage, uint32(len(queueItemData)))
	payloadMessage = append(payloadMessage, []byte(queueItemData)...)

	return payloadMessage
}

func sendMessage(t *testing.T, conn net.Conn, bytes []byte) {
	t.Helper()

	_, err := conn.Write(bytes)

	if err != nil {
		t.Fatalf("message writing error, %v", err)
	}

	_, err = io.ReadFull(conn, make([]byte, 1))

	if err != nil {
		t.Fatalf("read confirmation error, %v", err)
	}
}

func checkExpectedQueueItemsCount(t *testing.T, queueKey string, expectedQueueItemsCount int) {
	t.Helper()

	q, _ := queueManager.Get(queueKey)

	queueItemsCount, _ := q.Count()

	if queueItemsCount != expectedQueueItemsCount {
		t.Fatalf("the expected number of items in the queue is %v, %v items were received",
			expectedQueueItemsCount, queueItemsCount)
	}
}

func receiveQueueItem(t *testing.T, conn net.Conn) []byte {
	t.Helper()

	buf := make([]byte, 5)

	_, err := io.ReadFull(conn, buf)

	if err != nil {
		t.Fatalf("error reading header and data length, %v", err)
	}

	buf = make([]byte, binary.BigEndian.Uint32(buf[1:]))

	_, err = io.ReadFull(conn, buf)

	if err != nil {
		t.Fatalf("data read error, %v", err)
	}

	const confirmationType = 3

	_, err = conn.Write([]byte{confirmationType})

	if err != nil {
		t.Fatalf("confirmation write error, %v", err)
	}

	return buf
}
