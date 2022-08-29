# 🔋 AiRISTA Flow RTLS - BatteryHistory

BatteryHistory is a simple application that displays an interactive view of a battery discharge curve. Specifically, the battery history of any active RTLS tag can be viewed.

> The following images show the main views of the application.

<table style="border:none;">
  <tr>
    <td><img src="" width="480"/></td>
    <td><img src="" width="480"/></td>
  </tr>
</table>

## ⚙️ How it works

The application determines the data points for the curve by addressing the XML API of the RTLS controller. The MAC address of a tag is used as a filter. The API returns the entire battery history of the tag after authentication and request. The data is then plotted on a line graph.
## 🎯 Project goals

- [x] Display battery history of any tag
- [x] Provide data in web frontend
- [x] Interactive selection of the period to be viewed
- [ ] Screenshot function for the adjusted display
- [ ] Import all MAC addresses for easy selection

## 📜 Installation guide

The first thing to do is to clone the repository. After that you can simply run `go build` and start the compiled program in the terminal. At the first start the program creates a configuration file. Here the access data, port and IP address of the RTLS server must be entered. In addition, the port for the web server of the program must be selected.
