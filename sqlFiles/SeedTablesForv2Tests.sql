--clean up data from v1 tests or there will be problems
delete from page;
delete from pagesection;
delete from pagerevision;
delete from wikiuser where user_id > 4 ;

--add tag to allow duplicate test
select create_tag('oldtag','',2);
select create_tag('daltag','',2);