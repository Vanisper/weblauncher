; WebLauncher 安装程序脚本
; 由构建系统自动生成，不要手动编辑

#define MyAppId "{{APP_ID}}"
#define MyAppName "{{APP_NAME}}"
#define MyAppVersion "{{APP_VERSION}}"
#define MyAppPublisher "{{APP_PUBLISHER}}"
#define MyAppURL "{{APP_URL}}"
#define MyAppExeName "{{OUTPUT_NAME}}"

[Setup]
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}

DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes

OutputDir=..\..\.output
OutputBaseFilename={#MyAppName}_Setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern

PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

SetupIconFile=..\..\src\assets\icon.ico
UninstallDisplayIcon={app}\{#MyAppExeName}
UninstallDisplayName={#MyAppName}

[Languages]
Name: "chinesesimplified"; MessagesFile: "ChineseSimplified.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: checkedonce

[Files]
Source: "..\..\.output\{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\src\assets\icon.ico"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autoprograms}\{#MyAppName}\卸载 {#MyAppName}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

[Code]
function InitializeSetup(): Boolean;
begin
  if CheckForMutexes('WebLauncher_SingleInstance') then
  begin
    MsgBox('安装程序检测到 {#MyAppName} 正在运行。'#13#10'请先退出程序后再继续安装。', mbError, MB_OK);
    Result := false;
  end
  else
  begin
    Result := true;
  end;
end;

function InitializeUninstall(): Boolean;
begin
  if CheckForMutexes('WebLauncher_SingleInstance') then
  begin
    MsgBox('卸载程序检测到 {#MyAppName} 正在运行。'#13#10'请先退出程序后再继续卸载。', mbError, MB_OK);
    Result := false;
  end
  else
  begin
    Result := true;
  end;
end;
