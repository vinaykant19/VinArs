package database

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"runtime"
	cnf "../../../configuration"
	helper "../helper"
)
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)
type SetupDb struct {
	Db Database
}

type DbFields struct {
	Name string
	Type string //ast.Expr
	Length string
}

type DbTable struct {
	Name string
	Fields []DbFields
}

type DbSchema struct {
	DbTables []DbTable
}
func (dbSetup *SetupDb) Setup(configuration *cnf.Configuration, logger *log.Logger) (error){
	fmt.Println("DB Setup....")

	readDbEntity(configuration, logger)
	return nil
}

func readDbEntity(configuration *cnf.Configuration, logger *log.Logger) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, basepath + "/../../../src/entity", nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		logger.Println("Failed to parse package:", err)
	}

	DbSchemaData := DbSchema{
	}
	for _, pack := range packs {
		for _, file := range pack.Files {
			tableNames :=  file.Scope.Objects
			tableData := DbTable{}
			for _, obj := range tableNames {
				//fmt.Println(obj.Name) //Users
				tableData.Name =  obj.Name
				fields := obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
				fieldsData := DbFields{}
				for _, field := range fields {
					fmt.Println(field.Comment.Text())
					fieldType, fieldLength := helper.MapDbType(field.Type, "")
					fieldsData = DbFields{
						field.Names[0].String(),
						fieldType,
						fieldLength,
					}
					tableData.Fields = append(tableData.Fields, fieldsData)
					//fmt.Println(field.Names[0])
					//fmt.Println(field.Type)
				}
			}
			DbSchemaData.DbTables =  append(DbSchemaData.DbTables, tableData)
		}
	}

	fmt.Println(DbSchemaData)
}
