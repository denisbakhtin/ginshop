echo Commiting application

set /p name=Enter commit name: 

git add -u
git add .
git commit -m "%name%"
git push origin master
