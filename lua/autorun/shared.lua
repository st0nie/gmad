AddCSLuaFile()

local url = "http://your.server.domain:28080/"

if CLIENT then
    local function heartbeat()
        http.Fetch(url, function(body, len, headers, code)
            print("IDC: " .. body)
        end, function(error)
            print("Error: " .. error)
        end)
    end

    timer.Create("heartbeat", 30, 0, heartbeat)
end