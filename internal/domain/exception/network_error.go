package exception

import (
	"fmt"
	"net"
	"strings"
	"errors"
)

const (
	ErrorTypeNetworkUnavailable = "NETWORK_UNAVAILABLE"
	ErrorTypeDNSResolution      = "DNS_RESOLUTION_FAILED"
	ErrorTypeConnectionTimeout  = "CONNECTION_TIMEOUT"
	ErrorTypeConnectionRefused  = "CONNECTION_REFUSED"
	ErrorTypeTLSHandshake       = "TLS_HANDSHAKE_FAILED"
	ErrorTypeStorageService     = "STORAGE_SERVICE_UNAVAILABLE"
	ErrorTypeUploadFailed       = "UPLOAD_FAILED"
	ErrorTypeDownloadFailed     = "DOWNLOAD_FAILED"
	ErrorTypeThirdPartyService  = "THIRD_PARTY_SERVICE_UNAVAILABLE"
)

type NetworkError struct {
	Type        string
	Message     string
	Service     string // Which service failed (S3, API, etc.)
	Operation   string // What operation was being performed
	OriginalErr error
}

func (e NetworkError) Error() string {
	base := fmt.Sprintf("[%s] %s", e.Type, e.Message)
	
	if e.Service != "" {
		base = fmt.Sprintf("%s - Service: %s", base, e.Service)
	}
	
	if e.Operation != "" {
		base = fmt.Sprintf("%s - Operation: %s", base, e.Operation)
	}
	
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s (%s)", base, e.OriginalErr.Error())
	}
	
	return base
}

func (e NetworkError) Unwrap() error {
	return e.OriginalErr
}

// Generic network error
func NewNetworkError(message string, service string, operation string, originalErr error) *NetworkError {
	if message == "" {
		message = "Network operation failed"
	}
	
	return &NetworkError{
		Type:        ErrorTypeNetworkUnavailable,
		Message:     message,
		Service:     service,
		Operation:   operation,
		OriginalErr: originalErr,
	}
}

// DNS resolution error
func NewDNSResolutionError(host string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Failed to resolve DNS for host: %s", host)
	
	return &NetworkError{
		Type:        ErrorTypeDNSResolution,
		Message:     message,
		Service:     host,
		Operation:   "DNS Lookup",
		OriginalErr: originalErr,
	}
}

// Connection timeout error
func NewConnectionTimeoutError(service string, timeout string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Connection timeout to %s", service)
	if timeout != "" {
		message = fmt.Sprintf("%s after %s", message, timeout)
	}
	
	return &NetworkError{
		Type:        ErrorTypeConnectionTimeout,
		Message:     message,
		Service:     service,
		Operation:   "Connect",
		OriginalErr: originalErr,
	}
}

// Connection refused error
func NewConnectionRefusedError(service string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Connection refused by %s", service)
	
	return &NetworkError{
		Type:        ErrorTypeConnectionRefused,
		Message:     message,
		Service:     service,
		Operation:   "Connect",
		OriginalErr: originalErr,
	}
}

// TLS handshake error
func NewTLSHandshakeError(service string, originalErr error) *NetworkError {
	message := fmt.Sprintf("TLS handshake failed with %s", service)
	
	return &NetworkError{
		Type:        ErrorTypeTLSHandshake,
		Message:     message,
		Service:     service,
		Operation:   "TLS Handshake",
		OriginalErr: originalErr,
	}
}

// Storage service error (S3, MinIO, etc.)
func NewStorageServiceError(service string, bucket string, operation string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Storage service %s unavailable", service)
	if bucket != "" {
		message = fmt.Sprintf("%s - Bucket: %s", message, bucket)
	}
	
	return &NetworkError{
		Type:        ErrorTypeStorageService,
		Message:     message,
		Service:     service,
		Operation:   operation,
		OriginalErr: originalErr,
	}
}

// Upload failed error
func NewUploadFailedError(service string, fileName string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Failed to upload file to %s", service)
	if fileName != "" {
		message = fmt.Sprintf("%s: %s", message, fileName)
	}
	
	return &NetworkError{
		Type:        ErrorTypeUploadFailed,
		Message:     message,
		Service:     service,
		Operation:   "Upload",
		OriginalErr: originalErr,
	}
}

// Download failed error
func NewDownloadFailedError(service string, fileName string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Failed to download file from %s", service)
	if fileName != "" {
		message = fmt.Sprintf("%s: %s", message, fileName)
	}
	
	return &NetworkError{
		Type:        ErrorTypeDownloadFailed,
		Message:     message,
		Service:     service,
		Operation:   "Download",
		OriginalErr: originalErr,
	}
}

// Third-party service error
func NewThirdPartyServiceError(service string, endpoint string, originalErr error) *NetworkError {
	message := fmt.Sprintf("Third-party service %s unavailable", service)
	if endpoint != "" {
		message = fmt.Sprintf("%s - Endpoint: %s", message, endpoint)
	}
	
	return &NetworkError{
		Type:        ErrorTypeThirdPartyService,
		Message:     message,
		Service:     service,
		Operation:   "API Call",
		OriginalErr: originalErr,
	}
}

// Helper function to check if error is NetworkError
func IsNetworkError(err error) bool {
	_, ok := err.(*NetworkError)
	return ok
}

// Get network error type
func GetNetworkErrorType(err error) string {
	if netErr, ok := err.(*NetworkError); ok {
		return netErr.Type
	}
	return ""
}

// Parse and classify any error into NetworkError
func ClassifyNetworkError(err error, service string, operation string) *NetworkError {
	if err == nil {
		return nil
	}
	
	// If it's already our NetworkError, return it
	if netErr, ok := err.(*NetworkError); ok {
		return netErr
	}
	
	errStr := err.Error()
	
	// Check for DNS errors
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return NewDNSResolutionError(dnsErr.Name, err)
	}
	
	// Check for network timeout errors
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return NewConnectionTimeoutError(service, "", err)
	}
	
	// Check error strings
	switch {
	case strings.Contains(errStr, "no such host"):
		return NewDNSResolutionError(extractHost(errStr), err)
	case strings.Contains(errStr, "connection refused"):
		return NewConnectionRefusedError(service, err)
	case strings.Contains(errStr, "timeout"), strings.Contains(errStr, "deadline exceeded"):
		return NewConnectionTimeoutError(service, "", err)
	case strings.Contains(errStr, "TLS handshake"):
		return NewTLSHandshakeError(service, err)
	case strings.Contains(errStr, "server misbehaving"):
		return NewDNSResolutionError(extractHost(errStr), err)
	case strings.Contains(errStr, "bucket"), strings.Contains(errStr, "S3"), strings.Contains(errStr, "storage"):
		return NewStorageServiceError(service, extractBucket(errStr), operation, err)
	default:
		return NewNetworkError("Network operation failed", service, operation, err)
	}
}

// Helper to extract host from error message
func extractHost(errMsg string) string {
	// Try to extract host from common error patterns
	if strings.Contains(errMsg, "lookup") {
		parts := strings.Split(errMsg, "lookup ")
		if len(parts) > 1 {
			hostParts := strings.Split(parts[1], " ")
			if len(hostParts) > 0 {
				return hostParts[0]
			}
		}
	}
	return "unknown"
}

// Helper to extract bucket name from error message
func extractBucket(errMsg string) string {
	if strings.Contains(errMsg, "bucket") {
		parts := strings.Split(errMsg, "bucket ")
		if len(parts) > 1 {
			bucketParts := strings.Split(parts[1], " ")
			if len(bucketParts) > 0 {
				return bucketParts[0]
			}
		}
	}
	return ""
}