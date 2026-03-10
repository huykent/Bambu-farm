package camera

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"bambu-farm/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProxyService struct {
	logger         *zap.SugaredLogger
	printerService *service.PrinterService
	activeStreams  map[uint]string
	mu             sync.RWMutex
}

func NewProxyService(logger *zap.SugaredLogger, printerService *service.PrinterService) *ProxyService {
	return &ProxyService{
		logger:         logger,
		printerService: printerService,
		activeStreams:  make(map[uint]string),
	}
}

// StreamHandler returns an HTTP handler that proxies the video feed for the specified printer.
// For the sake of this prototype, we simulate an MJPEG stream or return placeholder content
// since full WebRTC/RTSP negotiation requires external binary dependencies (e.g. ffmpeg or go2rtc).
func (s *ProxyService) StreamHandler(c *gin.Context, orgID, printerID uint) {
	printer, err := s.printerService.GetPrinter(printerID, orgID)
	if err != nil || printer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Printer not found"})
		return
	}

	if printer.IPAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Printer IP not available"})
		return
	}

	rtspURL := fmt.Sprintf("rtsp://%s:8554/live", printer.IPAddress)
	s.logger.Infof("Initializing proxy stream for %s (URL: %s)", printer.Name, rtspURL)

	// Set headers for MJPEG streaming
	c.Writer.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Pragma", "no-cache")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	s.streamSimulatedMJPEG(ctx, c.Writer)
}

// streamSimulatedMJPEG provides a dummy stream loop. In a production system, this would:
// 1. Dial the RTSP URL via avformat or pure go RTSP library
// 2. Decode frames
// 3. Encode into MJPEG or wrap for WebRTC and push to the connection
func (s *ProxyService) streamSimulatedMJPEG(ctx context.Context, w io.Writer) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 1x1 black pixel JPEG magic bytes as a placeholder
	dummyJpeg := []byte{
		0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x48,
		0x00, 0x48, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0xff,
		0xc0, 0x00, 0x0b, 0x08, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xff, 0xc4, 0x00, 0x1f,
		0x00, 0x00, 0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0xff, 0xda, 0x00,
		0x08, 0x01, 0x01, 0x00, 0x00, 0x3f, 0x00, 0x37, 0xff, 0xd9,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Write the multipart boundary and header
			w.Write([]byte("--frame\r\n"))
			w.Write([]byte("Content-Type: image/jpeg\r\n\r\n"))
			w.Write(dummyJpeg)
			w.Write([]byte("\r\n"))
			
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}
