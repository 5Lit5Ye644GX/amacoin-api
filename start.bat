@echo off

set chain=test
set exe=amacoin-api.exe
set output=.env
set multi_dir=%APPDATA%\MultiChain\%chain%

echo [INFO] (T-7) The chain %chain% will be created if not existing in the %multi_dir% directory

if not exist %CD%\multichain goto download

:multichain
cd multichain
if not exist %multi_dir% goto initialize

:start
start multichaind.exe %chain%
echo [OK] (T-4) Multichain daemon has been launched on a new windows (don't close it)
cd ..

if not exist %output% goto configuration

:build
if not exist %exe% goto back
:launch
timeout 7
start %exe%
echo [FINISH] (T+1) The server is now launched! You can also launch it manually with ./%exe%
pause
exit

:download
echo [INFO] (T-6) Downloading MultiChain... Please wait...
powershell -Command "(New-Object Net.WebClient).DownloadFile('https://www.multichain.com/download/multichain-windows-2.0-alpha-2.zip', 'multichain.zip')"
powershell -Command "Expand-Archive multichain.zip"
del multichain.zip
echo [OK] Multichain successfully installed in ./multichain
goto multichain


:initialize
echo [INFO] (T-5) Initialize new Blockchain
multichain-util.exe create %chain%
echo [OK] The blockchain has been created and is ready for launch
goto start


:configuration
set conf=%multi_dir%\multichain.conf

for /f "delims== tokens=1,2" %%G in (%conf%) do set %%G=%%H
echo MULTICHAIN_CHAIN_NAME=%chain% >> %output%
echo MULTICHAIN_HOST=localhost >> %output%
echo MULTICHAIN_RPC_USER=%rpcuser% >> %output%
echo MULTICHAIN_RPC_PASSWORD=%rpcpassword% >> %output%

for /f "delims==  tokens=1,2" %%B in (%multi_dir%\params.dat) do set %%B=%%C
for /f "delims=#  tokens=1,2" %%E in ("%default-rpc-port %") do set port=%%E
echo MULTICHAIN_PORT=%port% >> %output%

echo [INFO] All Multichain parameters have been filled in %output% file
echo [OK] %output% is ready!
goto build

:back
echo [INFO] (T-1) Build server...
go build
echo [OK] You can now launch it with ./%exe%
goto launch