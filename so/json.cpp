#include "json.hpp"
#include <fstream>
#include <iostream>
#include <string>

using json = nlohmann::json;
json LoadJsonFile(std::string path) {
  std::ifstream f(path);
  json data = json::parse(f);
  return data;
}
json LoadJsonStr(std::string jsonStr) {

  json data = json::parse(jsonStr);

  return data;
}
int getJsonInt(json j, std::string key) {
  int num = j[key];
  return num;
}
std::string getJsonString(json j, std::string key) {
  std::string str = j[key];
  return str;
}
std::string *getJsonStringArray(json j, std::string key) {
  auto json_array = j[key];
  int json_length = json_array.size();
  //  std::string *strArray = new std::string[json_length];
  std::string *strArray = new std::string[json_length];
  for (int i = 0; i < json_length; i++) {
    strArray[i] = json_array[i];
  }
  return strArray;
}
int getStringArrayLength(std::string *array) {
  int len = sizeof(*array) / sizeof(std::string *);
  return len;
}
int getObjArrayLength(json j) {
  int len = j.size();
  return len;
}
json getJsonObj(json j, std::string key) {
  json jobj = j[key];
  return jobj;
}
bool getJsonBool(json j, std::string key) {
  bool b = j[key];
  return b;
}
