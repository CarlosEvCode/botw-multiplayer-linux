# The Legend of Zelda: Breath of the Wild Multiplayer para Linux

[English](README.md) | [Español](README.es.md)

Este repositorio provee un script de instalación automatizado y documentación técnica exhaustiva para ejecutar el mod multijugador (*Milk Bar Launcher* / *Breath of the Wild Multiplayer*) en distribuciones GNU/Linux y SteamOS bajo Wine.

---

## 1. Arquitectura del Sistema

El mod multijugador está compuesto por tres componentes principales:

1. **Emulador Cemu (v1.26.2 x86_64 para Windows)**:
   - Ejecuta la versión de Wii U del juego (*The Legend of Zelda: Breath of the Wild*).
   - Requiere la activación de Graphic Packs específicos (`BreathOfTheWild_BCML`, `bcmlPatches/MilkBarLauncher` y `ExtendedMemory`).

2. **Milk Bar Launcher (Cliente .NET / WinUI)**:
   - Controla el inicio de sesión y la conexión al servidor.
   - Inyecta la librería `InjectDLL.dll` en el proceso de Cemu y establece comunicación de lectura/escritura de memoria a través de tuberías con nombre (*Named Pipes*).

3. **MBL Dedicated Server (Servidor Dedicado .NET)**:
   - Procesa la sincronización de posiciones, animaciones, inventario, misiones y estados de juego entre clientes a través del puerto UDP/TCP `5050`.

---

## 2. Problemas de Compatibilidad en Linux y Soluciones Aplicadas

Durante el análisis y pruebas en entornos Linux, se diagnosticaron cinco puntos críticos de falla:

### A. Aislamiento de Procesos en Wine
* **Problema:** Si Cemu y Milk Bar Launcher se ejecutan en Wineprefixes separados o en instancias nativas independientes, Milk Bar Launcher no puede interactuar con la memoria ni con las tuberías con nombre de Cemu. Esto genera bloqueos indefinidos (*pantalla blanca "Cemu 1.26.2f - loading..."*).
* **Solución:** Ambos programas deben residir y ejecutarse estrictamente dentro del mismo Wineprefix de 64 bits (`~/.local/share/wineprefixes/botw-multiplayer`).

### B. Versión del Runtime de .NET
* **Problema:** Las versiones actuales de Milk Bar Launcher requieren `.NET Desktop Runtime 8.0 (x64)`. La presencia exclusiva de .NET 6.0 o versiones anteriores provoca el cierre de la aplicación con código de salida 150.
* **Solución:** Se instala `windowsdesktop-runtime-8.0.x-win-x64.exe` y `VC_redist.x64.exe` dentro del prefijo.

### C. Validación de Modelos en BCML ("Mod is not setup on BCML")
* **Problema:** Al presionar *Connect*, Milk Bar Launcher valida la existencia de 32 archivos de actores (`Jugador1.sbactorpack` a `Jugador32.sbactorpack`) dentro de la ruta `store_dir/merged/content/Actor/Pack/` definida en `settings.json`. Si dicha ruta no contiene los archivos, el cliente bloquea la conexión.
* **Solución:** Se genera el archivo de configuración `settings.json` con rutas de la unidad virtual `C:` y se enlaza la carpeta de Graphic Packs mergeados (`BreathOfTheWild_BCML`) a `AppData/Local/bcml/merged`.

### D. Crash del Servidor Dedicado por resolución de "localhost"
* **Problema:** En entornos Wine, la resolución de red de .NET puede fallar con `NullReferenceException` al enlazar con la dirección de texto `localhost`.
* **Solución:** En `ServerConfig.ini`, el parámetro `IP` debe fijarse explícitamente en `127.0.0.1`.

### E. Recursos faltantes en Roaming AppData
* **Problema:** El servidor dedicado intenta cargar `QuestFlagsNames.txt` y `ArmorMapping.txt` desde `%APPDATA%/BOTWM/`. Al extraerse con prefijos de ensamblado (`BOTWM.DedicatedServer.AppdataFiles.*`), el servidor lanza `FileNotFoundException`.
* **Solución:** Se copian y renombran automáticamente los archivos de configuración requeridos dentro de `AppData/Roaming/BOTWM/`.

---

## 3. Contenido del Repositorio y Releases

| Recurso | Ubicación | Descripción |
| :--- | :--- | :--- |
| `setup_botw_multiplayer.sh` | Repositorio | Script instalador automatizado para Linux. |
| `cemu_1.26.2.7z` | GitHub Releases | Distribución ligera (~146 MB) de Cemu 1.26.2 con Graphic Packs, 32 modelos de Link y parches BCML. |
| `MilkBarLauncher.zip` | GitHub Releases | Binarios del cliente Milk Bar Launcher, servidor dedicado y librerías de inyección. |
| `windowsdesktop-runtime-8.0.31-win-x64.exe` | GitHub Releases | Instalador oficial de Microsoft .NET Desktop Runtime 8.0 x64. |
| `VC_redist.x64.exe` | GitHub Releases | Instalador de Microsoft Visual C++ Redistributable 2015-2022 x64. |

---

## 4. Guía de Instalación

### Requisitos previos del sistema
Instalar las herramientas base mediante el gestor de paquetes de su distribución:

```bash
# En Arch / CachyOS / Manjaro:
sudo pacman -S wine winetricks p7zip curl

# En Ubuntu / Debian:
sudo apt install wine winetricks p7zip-full curl

# En Fedora:
sudo dnf install wine winetricks p7zip p7zip-plugins curl
```

### Instalación automatizada

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/CarlosEvCode/botw-multiplayer-linux.git
   cd botw-multiplayer-linux
   chmod +x setup_botw_multiplayer.sh
   ```

2. Ejecutar el script indicando la ruta absoluta de la carpeta de su juego base desempaquetado:
   ```bash
   ./setup_botw_multiplayer.sh "/ruta/absoluta/a/The Legend of Zelda Breath of the Wild"
   ```

*Nota: Si los archivos binarios no están presentes localmente, el script los descargará de forma automática desde la sección de Releases del repositorio.*

El script realiza de forma automatizada:
1. La inicialización del Wineprefix de 64 bits en `~/.local/share/wineprefixes/botw-multiplayer`.
2. La instalación silenciosa de Visual C++ y .NET 8.0 Desktop Runtime.
3. La extracción y estructura de Cemu 1.26.2 y Milk Bar Launcher en `drive_c`.
4. La vinculación del juego base, ajuste de archivos de configuración y enlaces simbólicos de BCML.
5. La creación de los scripts ejecutables en `~/Zelda_BotW_Multiplayer/` y accesos directos en el menú de aplicaciones (`~/.local/share/applications/`).

---

## 5. Guía de Uso y Conexión

Una vez completada la instalación, los scripts de control residen en `~/Zelda_BotW_Multiplayer/`:

1. **Iniciar Servidor Dedicado (solo el anfitrión / Host):**
   ```bash
   ~/Zelda_BotW_Multiplayer/1_iniciar_servidor.sh
   ```
   - Ingrese `0` para juego libre convencional o `1` para modos especiales (Prop Hunt, etc.).

2. **Iniciar Cliente (todos los jugadores):**
   ```bash
   ~/Zelda_BotW_Multiplayer/2_iniciar_milkbar.sh
   ```
   - **Host local:** Conectar a `127.0.0.1`, puerto `5050`.
   - **Clientes remotos (mediante VPN tipo Tailscale / ZeroTier):** Conectar a la IP del Host en dicha red, puerto `5050`.

3. **Configuración de Cemu (mandos / gráficos):**
   ```bash
   ~/Zelda_BotW_Multiplayer/3_iniciar_cemu.sh
   ```

---

## 6. Créditos y Atribución

- **Milk Bar Launcher & Dedicated Server**: Desarrollado y mantenido por la comunidad de [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher).
- **Cemu (Wii U Emulator)**: Desarrollado por [Team Cemu](https://cemu.info/) bajo licencia Mozilla Public License 2.0.
- **BCML (BotW Cross-Platform Mod Loader)**: Creado por [NiceneNerd](https://github.com/NiceneNerd/BCML).

### Descargo de Responsabilidad (Disclaimer)
Este repositorio y sus recursos asociados no contienen, alojan ni distribuyen archivos de juegos, ejecutables binarios oficiales (`.rpx`), llaves criptográficas propietarias ni contenidos protegidos por derechos de autor de Nintendo. Los usuarios deben proveer sus propios archivos legítimos del juego base, actualizaciones y DLCs.
