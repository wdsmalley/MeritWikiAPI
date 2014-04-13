CREATE TABLE Tag
	(
		Tag_ID serial,
		Tag_Name character(15),
		Tag_Description character(255),
		Created_ts timestamp with time zone,
		Created_by_user int,
		 CONSTRAINT tag_createdbyuser_fkey FOREIGN KEY (Created_by_user)
		  REFERENCES wikiuser (user_id) MATCH SIMPLE
		  ON UPDATE NO ACTION ON DELETE CASCADE	

	)
	;	
CREATE UNIQUE INDEX Tag_x01 ON Tag(Tag_Name)
	;
CREATE TYPE tagtype AS (tag_id int, 
						Tag_Name character(15),
						Tag_Description character(255));	
CREATE TABLE PageTag
	(
		Tag_ID int,
		Page_ID int,
		Created_ts timestamp with time zone,
		Created_by_user int,
		CONSTRAINT pagetag_pk PRIMARY KEY (tag_id,page_id),
		 CONSTRAINT pagetag_fkey FOREIGN KEY (Page_id)
		  REFERENCES page (pageid) MATCH SIMPLE
		  ON UPDATE NO ACTION ON DELETE CASCADE	

	)
	;		
	
CREATE TABLE UserTag
	(
		Tag_ID int,
		User_ID int,
		Created_ts timestamp with time zone,
		CONSTRAINT usertag_pk PRIMARY KEY (tag_id,user_id),
		 CONSTRAINT usertag_fkey FOREIGN KEY (user_id)
		  REFERENCES wikiuser (user_id) MATCH SIMPLE
		  ON UPDATE NO ACTION ON DELETE CASCADE	

	)
	;			
CREATE OR REPLACE FUNCTION create_wikiuser(iuser_name text, iuser_email text, iuser_pw text)
	RETURNS integer AS
	$BODY$
		DECLARE ouser_id integer;
		BEGIN
			if exists(select * 
			          from WikiUser
				      where user_email = iuser_email) 
            then
				ouser_id =  -1;
			else
				INSERT INTO WikiUser (
				  User_Name,
				  User_EMail ,
				  User_Password,
				  User_Status,
				  User_Rating,
				  created_ts,
				  last_update_ts) 
				  VALUES ( 
				  iUser_Name,
				  iUser_EMail,
				  iuser_PW,
				  0,
				  0,
				  current_timestamp,
				  current_timestamp
				);
				SELECT user_id INTO ouser_id
				FROM WikiUser
				WHERE User_email = iuser_email;
			end if;
			RETURN ouser_id;
		END
	$BODY$
	LANGUAGE plpgsql 
	;
CREATE OR REPLACE FUNCTION create_tag(itag_name text, itag_desc text, iuser_id int)
	RETURNS integer AS
	$BODY$
		DECLARE otag_id integer;
		BEGIN
			if exists(select * 
			          from Tag
				      where Tag_Name = itag_name) 
            then
				return -1;
			end if;
			if not exists(select * 
              		  from wikiuser
			          where user_id = iuser_id) 
            then
				return -1;
			end if;
				INSERT INTO Tag (
				  Tag_Name,
				  Tag_Description ,
				  created_ts,
				  Created_by_user) 
				  VALUES ( 
				  itag_name,
				  itag_desc,
				  current_timestamp,
				  iuser_id
				);
				SELECT tag_id INTO otag_id
				FROM Tag
				WHERE Tag_Name = itag_name;

			RETURN otag_id;
		END
	$BODY$
	LANGUAGE plpgsql 
	;
CREATE OR REPLACE FUNCTION delete_tag(itag_name text)
	RETURNS boolean AS
	$BODY$
		BEGIN
			if not exists(select * 
			          from Tag
				      where Tag_Name = itag_name) 
            then
				return false;
			end if;
			DELETE FROM Tag 
			WHERE Tag_Name = itag_name;

			RETURN true;
		END
	$BODY$
	LANGUAGE plpgsql 
	;
	
CREATE OR REPLACE FUNCTION create_usertag(itag_name text, iuser_id int)
	RETURNS boolean AS
	$BODY$
		DECLARE ltag_id integer;
		DECLARE oresult boolean;
		BEGIN
			if exists(select * 
			          from Tag
				      where Tag_Name = itag_name) 
			then
				SELECT tag_id INTO ltag_id
				FROM Tag
				WHERE Tag_Name = itag_name;
				INSERT INTO UserTag (
					Tag_ID,
					User_id ,
					created_ts) 
				VALUES ( 
					ltag_id,
					iuser_id,
					current_timestamp);
				oresult = true;
            else
				oresult = false;
			end if;
				
			RETURN oresult;
		END
	$BODY$
	LANGUAGE plpgsql 
	;
CREATE OR REPLACE FUNCTION create_pagetag(itag_name text, ipage_id int, iuser_id int)
	RETURNS boolean AS
	$BODY$
		DECLARE ltag_id integer;
		DECLARE oresult boolean;
		BEGIN
			if exists(select * 
			          from Tag
				      where Tag_Name = itag_name) 
			then
				SELECT tag_id INTO ltag_id
				FROM Tag
				WHERE Tag_Name = itag_name;
				INSERT INTO PageTag (
					Tag_ID,
					page_id ,
					created_ts,
					Created_by_user) 
				VALUES ( 
					ltag_id,
					ipage_id,
					current_timestamp,
					iuser_id);
				oresult = true;
            else
				oresult = false;
			end if;
				
			RETURN oresult;
		END
	$BODY$
	LANGUAGE plpgsql 
	;	
CREATE OR REPLACE FUNCTION get_all_tags()
	RETURNS setof tagtype AS
	$BODY$
		BEGIN
			return query SELECT 
				TAG_ID,
				TAG_NAME,
				TAG_DESCRIPTION
			FROM TAG;
		END
	$BODY$
	LANGUAGE plpgsql 
	;		
--delete existing user
DROP FUNCTION delete_user(text);
CREATE OR REPLACE FUNCTION delete_user(iUserEmail text)
  RETURNS boolean AS
$BODY$
	BEGIN
		if exists(
			SELECT * FROM Wikiuser
			WHERE User_Email = iUserEmail)
		then
			DELETE FROM Wikiuser
			WHERE User_Email = iUserEmail;
		else
			return false;
		end if;
		return true;
	END
$BODY$
  LANGUAGE plpgsql ;	
UPDATE AppVersion SET Version = 2 WHERE Name = 'wikidb'	;