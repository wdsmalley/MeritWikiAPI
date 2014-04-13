	CREATE TABLE AppVersion
	(Name character(30),
	 Version integer)
	;

	--create user table
	CREATE TABLE WikiUser
	(
		User_Name character(30),
		User_ID serial,
		User_EMail character(50),
		User_Password character(16),
		User_Status integer,
		User_rating integer,
		Last_logon_ts timestamp with time zone,
		Last_logoff_ts timestamp with time zone,
		Created_ts timestamp with time zone,
		Last_Update_ts timestamp with time zone
	)
	;
	CREATE UNIQUE INDEX wikiuser_x01 ON wikiuser(user_email)
	;
	ALTER TABLE wikiuser
	ADD CONSTRAINT wikiuser_pk PRIMARY KEY(user_id)
	;
	  --create page table
	CREATE TABLE page
	(
	  pageid serial NOT NULL,
	  pagetitle character(50),
	  pageurl character(50),
	  currentrevision integer,
	  createdbyuser integer,
	  created_ts timestamp with time zone,
	  lastupdateuser integer,
	  lastupdate_ts timestamp with time zone,  
	  CONSTRAINT page_pk PRIMARY KEY (pageid),
	  CONSTRAINT page_createdbyuser_fkey FOREIGN KEY (createdbyuser)
		  REFERENCES wikiuser (user_id) MATCH SIMPLE
		  ON UPDATE NO ACTION ON DELETE NO ACTION
	)
	;
	  -------------------------------------
	--create PageRevision table
	CREATE TABLE PageRevision
	(
	  PageID integer,
	  RevisionID integer,
	  RevisionCreatedByUser integer,
	  Revision_Created_ts timestamp with time zone
	)
	;

	--create PageSection table
	CREATE TABLE PageSection
	(
	  PageID integer,
	  RevisionID integer,
	  SectionID integer,
	  SectionTitle text,
	  SectionText text,
	  SectionCreatedByUserID integer,
	  SectionCreated_ts timestamp with time zone,
	  SectionLastUpdateByUserID integer,
	  SectionLastUpdate_ts timestamp with time zone
	)
	;
	--create the page view
	CREATE VIEW PageView AS
    SELECT
        P.pageid,
        P.pagetitle,
		P.pageurl,
		P.currentrevision,
		w.user_name,
		P.created_ts 
    FROM Page P 
    JOIN Wikiuser w 
	ON P.createdbyuser = w.user_id
    ;
    CREATE TYPE pageviewtype AS (pageid int, 
								 pagetitle character(50),
								 pageurl character(50), 
								 currentrevision int, 
								 user_name character(30), 
								 created_ts timestamp with time zone);
	CREATE TYPE usertype AS (user_id int, 
								 user_name character(30),
								 user_email character(50), 
								 user_status int, 
								 user_rating int);
	      ---////////////////////////////////
      --update the version table so this doesnt run again
      update AppVersion
      set Version = 1
      where Name = 'MeritWiki';
    
	  --////
	  --create functions
	  --add new user
	CREATE OR REPLACE FUNCTION create_wikiuser(iuser_name text, iuser_email text, iuser_pw text)
	RETURNS integer AS
	$BODY$
		DECLARE ouser_id integer;
		BEGIN
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
			RETURN ouser_id;
		END
	$BODY$
	LANGUAGE plpgsql 
	;
	


	
	--update rating of existing user
	CREATE OR REPLACE FUNCTION update_userRating(iuser_id integer, iuser_rating integer)
    RETURNS void AS
	$BODY$
		BEGIN
			UPDATE WikiUser 
			SET User_rating = iuser_Rating, last_update_ts = current_timestamp
			WHERE User_id = iUser_ID;
		END
	$BODY$
	LANGUAGE plpgsql;
	--update status of existing user
	CREATE OR REPLACE FUNCTION update_userStatus(iuser_id integer, iuser_status integer)
    RETURNS void AS
	$BODY$
		BEGIN
			UPDATE WikiUser 
			SET User_status = iuser_Status, last_update_ts = current_timestamp
			WHERE User_id = iUser_ID;
		END
	$BODY$
	LANGUAGE plpgsql;
	--logon user
	CREATE OR REPLACE FUNCTION logon_user(iUser_Email text,iUserPassword text)
    RETURNS boolean AS
	$BODY$
		BEGIN
			if not exists(select * 
			              from WikiUser
                          where User_Email= iUser_Email
                          and User_password = iUserPassword) then
				RETURN false;
			end if;
			UPDATE WikiUser 
			SET user_status = 1, last_logon_ts = current_timestamp,last_logoff_ts = null
			WHERE User_Email= iUser_Email
			and User_password = iUserPassword;
			RETURN true;
		END
	$BODY$
	LANGUAGE plpgsql;
	--logoff user
		CREATE OR REPLACE FUNCTION logoff_user(iuser_id integer)
    RETURNS void AS
	$BODY$
		BEGIN
			UPDATE WikiUser 
			SET user_status = 0, last_logoff_ts = current_timestamp
			WHERE User_id = iUser_ID;
		END
	$BODY$
	LANGUAGE plpgsql;
	

--lookup the user name by the input email	
CREATE OR REPLACE FUNCTION lookup_username(iuser_email text)
  RETURNS text AS
$BODY$
	DECLARE ouser_name text;
		BEGIN
			select user_name into ouser_name from wikiuser where user_email = iuser_email;
			return ouser_name;
		END
		
	$BODY$
  LANGUAGE plpgsql ;

  --lookup the user by email
CREATE OR REPLACE FUNCTION lookup_user(iuser_email text)
  RETURNS SETOF usertype AS
$BODY$
    BEGIN
	RETURN QUERY
        SELECT
            user_id, 
            user_name,
	        user_email,
	        user_status,
	        user_rating
	FROM WikiUser
	WHERE user_email = iuser_email;
    END
$BODY$
LANGUAGE plpgsql ;
--delete existing user
CREATE OR REPLACE FUNCTION delete_user(iUserName text)
  RETURNS void AS
$BODY$
	BEGIN
		DELETE FROM Wikiuser
		WHERE User_name = iUserName;
	END
$BODY$
  LANGUAGE plpgsql ;
  
--create a new page
CREATE OR REPLACE FUNCTION create_page(ipage_title text, ipage_url text, iuser_id integer)
  RETURNS integer AS
$BODY$
	DECLARE oPageID integer;
	BEGIN
	if not exists(select * 
	              from Wikiuser
			      where User_id = iuser_id) 
    then
		oPageID = 0;
	else
		INSERT INTO Page(
			pagetitle ,
			pageurl ,
			currentrevision,
			createdbyuser,
			created_ts 
		) 
		VALUES ( 
			ipage_title,
			ipage_url,
			1,
			iuser_id,
			current_timestamp
		);
		select pageid into oPageID
		from Page
		where pagetitle = ipage_title;
	end if;
	RETURN oPageID;
	END
$BODY$
  LANGUAGE plpgsql ;
--retrieve page info by title
CREATE OR REPLACE FUNCTION lookup_page_bytitle(iTitle text)
    RETURNS SETOF pageviewtype
    AS $$
    BEGIN
    
        RETURN QUERY
        SELECT 
            pageid,
            pagetitle,
	        pageurl,
	        currentrevision,
	        user_name,
	        created_ts 
	    FROM pageview
	    WHERE pagetitle = ititle;
    END;
    $$  
language plpgsql;

-- delete a page

CREATE OR REPLACE FUNCTION delete_page(iPageID integer)
  RETURNS void AS
$BODY$
	BEGIN
		DELETE FROM Page
		WHERE Pageid = iPageID;
	END
$BODY$
  LANGUAGE plpgsql ;

--create page sections

CREATE OR REPLACE FUNCTION Insert_PageSection(iPage_id integer,iSectionum integer, iSection_title text, iSection_text text,iUser_id integer)
    RETURNS integer AS $$
		DECLARE newRevisionID integer;
		BEGIN
			if not exists(select * 
			              from PageSection
				      where Pageid = iPage_id) 
                        then
				newRevisionID = 1;
			else 
				newRevisionID = (select max(RevisionID) + 1 from PageSection where Pageid = iPage_id);
			end if;
			INSERT INTO PageSection (pageid ,
			                         revisionid ,
			                         sectionid,
			                         sectiontitle ,
			                         sectiontext ,
			                         sectioncreatedbyuserid,
			                         sectioncreated_ts)
			VALUES(
				iPage_id,
				newRevisionID,
				iSectionum,
				iSection_title,
				iSection_text ,
				iUser_id,
				current_timestamp
			);
			SELECT revisionid INTO newRevisionID 
			FROM PageSection
			WHERE PageID = iPage_id AND
				SectionID = iSectionum AND
				RevisionID = newRevisionID;
			RETURN newRevisionID;
		END
	$$
	LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION update_page(iPage_id integer,iTitle text,iUser_id integer)
    RETURNS void AS $$
	BEGIN
    	    UPDATE Page
    	    SET PageTitle =iTitle,
				lastupdateuser = iUser_id,
				lastupdate_ts = current_timestamp
			WHERE PageID = iPage_id;
	END
$$
LANGUAGE plpgsql;	
CREATE OR REPLACE FUNCTION update_section(iPage_id integer,isectionid integer,iTitle text,itext text, iUser_id integer)
    RETURNS void AS $$
	DECLARE newRevisionID integer;
	BEGIN
			if not exists(select * 
			              from PageSection
				          where Pageid = iPage_id) 
            then
				newRevisionID = 1;
			else 
				newRevisionID = (select max(RevisionID) + 1 from PageSection where Pageid = iPage_id);
			end if;
    	    UPDATE PageSection
    	    SET SectionTitle =iTitle,
				SectionText =iText,
				RevisionID = newRevisionID,
				SectionLastUpdateByUserID = iUser_id,
				SectionLastUpdate_ts = current_timestamp
			WHERE PageId = iPage_ID and sectionid = isectionid ;
	END
$$
LANGUAGE plpgsql;	
--change password of existing user
	CREATE OR REPLACE FUNCTION update_wikiuser(iuser_id integer, iuser_name text, iemail text, iuser_pw text)
    RETURNS void AS $$
	DECLARE newUserName text;
	DECLARE newEmail text;
	DECLARE newPW text;

		BEGIN
			if iuser_pw != ''
			then
				newPW = iuser_pw;
			else
				newPW = (select User_Password from WikiUser WHERE User_id = iUser_ID);
			end if;
			
			if iemail != ''
			then
				newEmail = iemail;
			else
				newEmail = (select User_EMail from WikiUser WHERE User_id = iUser_ID);
			end if;
			
			if iuser_name != ''
			then
				newUserName = iuser_name;
			else
				newUserName = (select User_Name from WikiUser WHERE User_id = iUser_ID);
			end if;
			
			UPDATE WikiUser 
			SET User_Password = newPW,
				User_Name = newUserName,
				User_EMail = newEmail,
				last_update_ts = current_timestamp
			WHERE User_id = iUser_ID;
		END
	$$
	LANGUAGE plpgsql;
INSERT INTO APPVERSION(name,version) VALUES('wikidb',1);	