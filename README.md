This is a GoLang Template that will be used 
as FDS ASYA PHILIPPINES INC. standard format.

#******************************************#
#*********** GoLang is Required ***********#
#******************************************#

1.  Download the template by using either: git clone https://github.com/FDSAP-Git-Org/Go_Template.git or download the repository at Github
2.  Open go.mod file and change Template into desired project name
3.  Repeat step 2 in Makefile, change Template into desired project name
4.  Check other files that uses Template ang change it to desired project name
5.  Export files in folder project_env_files to root folder
6.  Open terminal and run the project by using make commad as listed below:
make DEV    (to run project in Development)
make SIT    (to run project in System Integration Testing)
make UAT    (to run project in User Acceptance Testing)
make PROD   (to run project in Production)
make KILLS  (to kill the running project)
make CHECK  (to check the port where the project is running)
make LOG    (to check or tail the logs of the project)