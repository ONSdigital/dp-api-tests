CREATE (node:`_code_list`:`_name_test_sex`:`_api_test` { label:'sex', edition:'one-off' });
MATCH (parent:`_code_list`:`_name_test-sex`) WITH parent CREATE (node:`_code`:`_api_test` { value:'2' })-[:usedBy { label:"Female"}]->(parent);
MATCH (parent:`_code_list`:`_name_test-sex`) WITH parent CREATE (node:`_code`:`_api_test` { value:'1' })-[:usedBy { label:"Male"}]->(parent);
MATCH (parent:`_code_list`:`_name_test-sex`) WITH parent CREATE (node:`_code`:`_api_test` { value:'0' })-[:usedBy { label:"All"}]->(parent);