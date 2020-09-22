package infocmdb

import (
	"errors"
	"strconv"
	"strings"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilError "github.com/infonova/infocmdb-sdk-go/util/error"

	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type getCiTypeIdByCiTypeName struct {
	Data []responseId `json:"data"`
}

func (c *Client) GetCiTypeIdByCiTypeName(name string) (r int, err error) {
	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetCiTypeIdByCiTypeName_" + name
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(int), nil
	}

	params := map[string]string{
		"argv1": name,
	}

	response := getCiTypeIdByCiTypeName{}
	err = c.v2.Query("int_getCiTypeIdByCiTypeName", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(name + " - " + v2.ErrNoResult.Error())
	case 1:
		r = response.Data[0].Id
		c.v1.Cache.Set(cacheKey, r, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(name + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type getCiTypeName struct {
	Data []struct {
		Name string `json:"name"`
	} `json:"data"`
}

func (c *Client) GetCiTypeName(ciId int) (ciTypeName string, err error) {
	ciIdString := strconv.Itoa(ciId)

	if err = c.v2.Login(); err != nil {
		return
	}

	cacheKey := "GetCiTypeName_" + ciIdString
	cached, found := c.v1.Cache.Get(cacheKey)
	if found {
		return cached.(string), nil
	}

	params := map[string]string{
		"argv1": ciIdString,
		"argv2": "name",
	}
	response := getCiTypeName{}
	err = c.v2.Query("int_getCiTypeOfCi", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	switch len(response.Data) {
	case 0:
		err = utilError.FunctionError(ciIdString + " - " + v2.ErrNoResult.Error())
	case 1:
		ciTypeName = response.Data[0].Name
		c.v1.Cache.Set(cacheKey, ciTypeName, utilCache.DefaultExpiration)
	default:
		err = utilError.FunctionError(ciIdString + " - " + v2.ErrTooManyResults.Error())
	}

	return
}

type respSetTypeOfCi struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *Client) SetTypeOfCi(ciId int, ciType string) (err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	currentCiTYpe, err := c.GetCiTypeName(ciId)
	if err != nil {
		return err
	}

	if ciType == currentCiTYpe {
		return errors.New("the requested ci type is already set: " + currentCiTYpe)
	}

	ciTypeId, err := c.GetCiTypeIdByCiTypeName(ciType)
	if err != nil {
		return err
	}

	ciTypeIdString := strconv.Itoa(ciTypeId)
	ciIdString := strconv.Itoa(ciId)

	params := map[string]string{
		"argv1": ciIdString,
		"argv2": ciTypeIdString,
		"argv3": "0",
	}

	response := respSetTypeOfCi{}
	err = c.v2.Query("int_setCiTypeOfCi", &response, params)
	if err != nil {
		err = utilError.FunctionError(err.Error())
		log.Error("Error: ", err)
		return
	}

	if response.Success != true {
		return errors.New("couldn't change ci type to: " + ciType + " for ciid: " + ciIdString + " ,error: " + response.Message)
	}

	return
}

type respCreateCiType struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []responseId `json:"data"`
}

type ciTypeParams struct {
	Name                    string
	Description             string
	Note                    string
	ParentCiTypeId          int
	OrderNumber             int
	CreateButtonDescription string
	Icon                    string
	Query                   string
	DefaultProjectId        int
	DefaultAttributeId      int
	DefaultSortAttributeId  int
	IsDefaultSortAsc        int
	IsCiAttach              int
	IsAttributeAttach       int
	Tag                     string
	IsTabEnabled            int
	IsEventEnabled          int
	IsActive                int
	UserId                  int
}

func (c *Client) NewTypeParams(options ...func(*ciTypeParams)) *ciTypeParams {
	typeParams := &ciTypeParams{
		Name:                    "",
		Description:             "",
		Note:                    "",
		ParentCiTypeId:          0,
		OrderNumber:             0,
		CreateButtonDescription: "",
		Icon:                    "",
		Query:                   "",
		DefaultProjectId:        0,
		DefaultAttributeId:      0,
		DefaultSortAttributeId:  0,
		IsDefaultSortAsc:        0,
		IsCiAttach:              0,
		IsAttributeAttach:       0,
		Tag:                     "",
		IsTabEnabled:            0,
		IsEventEnabled:          0,
		IsActive:                1,
		UserId:                  0,
	}

	for _, option := range options {
		option(typeParams)
	}
	return typeParams

}

func (c *Client) TypeWithName(name string) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.Name = name
	}
}

func (c *Client) TypeWithDesc(desc string) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.Description = desc
	}
}

func (c *Client) TypeWithNote(note string) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.Note = note
	}
}

func (c *Client) TypeWithParentCiTypeId(parentId int) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.ParentCiTypeId = parentId
	}
}

func (c *Client) TypeWithOrderNumber(orderNum int) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.OrderNumber = orderNum
	}
}

func (c *Client) TypeWithDefProjectId(projId int) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		obj.DefaultProjectId = projId
	}
}

func (c *Client) TypeIsActive(active bool) func(*ciTypeParams) {
	return func(obj *ciTypeParams) {
		if active != true {
			obj.IsActive = 0
		}
	}
}

func (c *Client) CreateCiType(typeParams *ciTypeParams) (typeId int, err error) {

	if err = c.v2.Login(); err != nil {
		return
	}

	typeIdExists, _ := c.GetCiTypeIdByCiTypeName(typeParams.Name)

	if typeIdExists == 0 {

		columns := []string{
			"name",
			"description",
			"note",
			"parent_ci_type_id",
			"order_number",
			"create_button_description",
			"icon",
			"query",
			"default_project_id",
			"default_attribute_id",
			"default_sort_attribute_id",
			"is_default_sort_asc",
			"is_ci_attach",
			"is_attribute_attach",
			"tag",
			"is_tab_enabled",
			"is_event_enabled",
			"is_active",
			"user_id",
		}

		values := []string{
			typeParams.Name,
			typeParams.Description,
			typeParams.Note,
			strconv.Itoa(typeParams.ParentCiTypeId),
			strconv.Itoa(typeParams.OrderNumber),
			typeParams.CreateButtonDescription,
			typeParams.Icon,
			typeParams.Query,
			strconv.Itoa(typeParams.DefaultProjectId),
			strconv.Itoa(typeParams.DefaultAttributeId),
			strconv.Itoa(typeParams.DefaultSortAttributeId),
			strconv.Itoa(typeParams.IsDefaultSortAsc),
			strconv.Itoa(typeParams.IsCiAttach),
			strconv.Itoa(typeParams.IsAttributeAttach),
			typeParams.Tag,
			strconv.Itoa(typeParams.IsTabEnabled),
			strconv.Itoa(typeParams.IsEventEnabled),
			strconv.Itoa(typeParams.IsActive),
			strconv.Itoa(typeParams.UserId),
		}

		params := map[string]string{
			"argv1": "`" + strings.Join(columns, "`, `") + "`",
			"argv2": "'" + strings.Join(values, "', '") + "'",
		}

		response := respCreateCiType{}
		err = c.v2.Query("int_createCIType", &response, params)
		if err != nil {
			err = utilError.FunctionError(err.Error())
			log.Error("Error: ", err)
			return
		}

		if response.Success != true {
			return 0, errors.New("error: " + response.Message)
		}

		switch len(response.Data) {
		case 0:
			err = utilError.FunctionError(typeParams.Name + " - " + v2.ErrNoResult.Error())
		case 1:
			typeId = response.Data[0].Id
		default:
			err = utilError.FunctionError(typeParams.Name + " - " + v2.ErrTooManyResults.Error())
		}

	} else {
		return typeIdExists, nil
	}

	return
}
