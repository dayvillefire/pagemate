# SITE DOCUMENTATION

Starting from the main URL:

* `$BASEURL`
  - form name="login" target="login.asp" method="post"
  - type="text" name="loginSubscriber" (UPPERCASE)
  - type="password" name="password" 
  - type="submit" name="submit" value="Submit"
* `$BASEURL/message/find_recipients.asp`
  - form name="f" action="list_recipients.asp" method="post"
  - type="text" name="objectname" id="ext-gen8"
  - type="submit" name="submit" value="Submit"
  - listing:
    - table class="labels" -> tbody -> #4...
    - tr -> | td -> a href="" | td -> &nbsp; | td -> a href=""
* `$BASEURL/message/send.asp?objectname=$GROUPNAME`
  - form name="f" target="do_send.asp" method="post"
  - type="text" name="display_objectname" value="$GROUPNAME"
  - type="text" name="description" value="$GROUPDESC"
  - type="textarea" name="comments"
  - type="text" name="message" id="ext-gen8"
  - type="text" name="reference"
  - type="submit" name="submit" value="Send"
