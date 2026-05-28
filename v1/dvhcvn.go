package dvhcvn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Level3 represents a ward (phường/xã)
type Level3 struct {
	Level3ID string `json:"level3_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// Level2 represents a district (quận/huyện)
type Level2 struct {
	Level2ID string    `json:"level2_id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Level3s  []Level3  `json:"level3s"`
}

// Level1 represents a province (tỉnh/thành phố)
type Level1 struct {
	Level1ID string    `json:"level1_id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Level2s  []Level2  `json:"level2s"`
}

type Service struct {
	remotePath string
}

func NewService(remotePath string) *Service {
	return &Service{remotePath: remotePath}
}

// GetProvinces returns all provinces from the remote data
func (s *Service) GetProvinces() ([]Level1, error) {
	data, err := s.fetchData()
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// GetProvince returns a specific province by ID
func (s *Service) GetProvince(provinceID string) (*Level1, error) {
	provinces, err := s.GetProvinces()
	if err != nil {
		return nil, err
	}
	
	for _, province := range provinces {
		if province.Level1ID == provinceID {
			return &province, nil
		}
	}
	
	return nil, fmt.Errorf("province with ID %s not found", provinceID)
}

// GetDistricts returns all districts from a specific province
func (s *Service) GetDistricts(provinceID string) ([]Level2, error) {
	province, err := s.GetProvince(provinceID)
	if err != nil {
		return nil, err
	}
	
	return province.Level2s, nil
}

// GetDistrict returns a specific district by ID from a province
func (s *Service) GetDistrict(provinceID, districtID string) (*Level2, error) {
	districts, err := s.GetDistricts(provinceID)
	if err != nil {
		return nil, err
	}
	
	for _, district := range districts {
		if district.Level2ID == districtID {
			return &district, nil
		}
	}
	
	return nil, fmt.Errorf("district with ID %s not found in province %s", districtID, provinceID)
}

// GetWards returns all wards from a specific district
func (s *Service) GetWards(provinceID, districtID string) ([]Level3, error) {
	district, err := s.GetDistrict(provinceID, districtID)
	if err != nil {
		return nil, err
	}
	
	return district.Level3s, nil
}

// GetWard returns a specific ward by ID from a district
func (s *Service) GetWard(provinceID, districtID, wardID string) (*Level3, error) {
	wards, err := s.GetWards(provinceID, districtID)
	if err != nil {
		return nil, err
	}
	
	for _, ward := range wards {
		if ward.Level3ID == wardID {
			return &ward, nil
		}
	}
	
	return nil, fmt.Errorf("ward with ID %s not found in district %s, province %s", wardID, districtID, provinceID)
}

// DVHCResponse represents the response structure from the API
type DVHCResponse struct {
	Data []Level1 `json:"data"`
}

// fetchData fetches and parses JSON data from the remote path
func (s *Service) fetchData() ([]Level1, error) {
	resp, err := http.Get(s.remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var response DVHCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	return response.Data, nil
}
