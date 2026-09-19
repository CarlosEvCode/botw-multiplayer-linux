# The Legend of Zelda: Breath of the Wild Multiplayer en Linux / SteamOS

[English (README.md)](README.md) | [Español (README.es.md)](README.es.md)

Instalador automatizado, TUI Manager interactivo y entorno técnico completo para hostear y jugar a *The Legend of Zelda: Breath of the Wild Multiplayer* (*Milk Bar Launcher* / *Cemu 1.26.2*) en distribuciones GNU/Linux y SteamOS mediante Wine.

---

## 1. Instalación Rápida

Ejecuta este comando en tu terminal para descargar, configurar e instalar todo automáticamente:

```bash
curl -sSL https://raw.githubusercontent.com/CarlosEvCode/botw-multiplayer-linux/main/install.sh | bash
```

### Instalación Manual (Vía Git)

```bash
git clone https://github.com/CarlosEvCode/botw-multiplayer-linux.git
cd botw-multiplayer-linux
chmod +x install.sh
./install.sh
```

> **Nota:** El instalador aprovisiona automáticamente Wine de 64 bits, Microsoft .NET 8.0 Desktop Runtime, Visual C++ 2015–2022, Cemu 1.26.2 (con los parches BCML y los 32 modelos de Link ya fusionados), Milk Bar Launcher (Servidor y Cliente Headless) y la interfaz CLI `botw-manager`.

---

## 2. TUI Manager (`botw-manager`)

Administra servidores, clientes, modos de juego, rutas del juego y direcciones IP de red directamente desde tu terminal:

```bash
botw-manager
```

### Funciones y Atajos de Teclado
* **Navegación Global:** Presiona `Tab` / `Shift+Tab` o las teclas numéricas `[1]`–`[6]` para alternar entre pestañas en cualquier momento.
* **1. Dashboard:**
  * `[s]`: Iniciar / Detener Servidor Dedicado (con consola de logs en vivo).
  * `[c]`: Conectar al Juego (inicia Cemu e inyecta el mod multijugador de forma headless mediante `MilkBar.CLI`).
  * `[x]`: Desconectar / Detener la sesión activa del cliente.
  * `[e]`: Lanzar Cemu 1.26.2 (Individual / Sin conexión).
  * `[t]`: Copiar IP de Tailscale / ZeroTier / LAN al portapapeles.
  * `[q]`: Salir de la aplicación.
* **2. Cliente & Conexión:**
  * Configura los parámetros de conexión multijugador: IP del servidor destino, puerto, contraseña, nombre del jugador y modelo del personaje.
  * `[c]` / `[Enter]`: Conectar al Juego.
  * `[l]`: Fijar IP destino a `127.0.0.1` (Anfitrión Local).
  * `[t]`: Fijar IP destino a la IP detectada de Tailscale/ZeroTier/LAN.
  * `[s]`: Guardar parámetros de conexión.
  * `[x]`: Desconectar cliente activo.
* **3. Modos de Juego:**
  * `←` / `→`: Alternar modos de juego:
    * **Cooperativo Estándar (Libre):** Modo cooperativo con sincronizaciones totalmente personalizables.
    * **Hunter vs Speedrunner:** Modo competitivo Manhunt (bloquea el progreso individual automáticamente).
    * **DeathSwap:** Modo de supervivencia con intercambio de posiciones periódicas (bloquea sincronizaciones de progreso).
  * `Espacio`: Alternar reglas de sincronización individuales (Misiones, Santuarios, Torres, Kologs, Enemigos, Mazmorras, Ubicaciones).
  * `s`: Guardar reglas (reinicia el servidor automáticamente si está corriendo).
* **4. Rutas & DLC:**
  * `↓` / `↑`: Seleccionar campo (Juego Base, Update v208, DLC v80).
  * `f`: Abrir el explorador de carpetas gráfico del sistema (*Zenity* / *Kdialog*) para elegir rutas con un clic.
  * `Enter`: Guardar y reconstruir los enlaces simbólicos de Wine y BCML.
* **5. Red:**
  * Muestra las interfaces activas de Tailscale VPN, Red Local LAN y ZeroTier con instrucciones de conexión.
* **6. Ajustes:**
  * **Selector de Idioma:** Cambia instantáneamente entre **English** y **Español**.
  * Información del proyecto, versión y enlaces al repositorio.

---

## 3. Guía de Conexión Multijugador

### A. Anfitrión Local (Jugando en la misma PC)
1. Abre `botw-manager` y presiona `[s]` para iniciar el Servidor Dedicado.
2. Presiona `[c]` para conectar inmediatamente (o ve a la Pestaña `2. Cliente & Conexión` y verifica `127.0.0.1:5050`).
3. Cemu se abrirá automáticamente sincronizado sin ninguna ventana gráfica intermedia de Wine.

### B. Red Local (LAN / Misma Wi-Fi o Ethernet)
1. El Anfitrión inicia el Servidor Dedicado en `botw-manager` (presionando `[s]`).
2. El Anfitrión comparte su **IP Local LAN** mostrada en `botw-manager` (ej. `192.168.1.50`).
3. Los demás jugadores abren `botw-manager` en sus equipos, van a la Pestaña `2. Cliente & Conexión`, ingresan la IP LAN del Anfitrión (`192.168.1.50`) y presionan `[c]`.

### C. Internet / Amigos Remotos (VPN Tailscale o ZeroTier)
1. Conéctate a la misma red virtual ([Tailscale](https://tailscale.com/) o [ZeroTier](https://www.zerotier.com/)) tanto en el equipo del Anfitrión como en el de los clientes.
2. El Anfitrión inicia el Servidor Dedicado y copia su IP de VPN (presionando `[t]` en `botw-manager`).
3. Los jugadores remotos abren `botw-manager`, van a la Pestaña `2. Cliente & Conexión`, pegan la IP de VPN del Anfitrión y presionan `[c]`.

---

## 4. Arquitectura del Cliente Headless y Estabilidad en Wayland

El lanzador original de Milk Bar para Windows utilizaba una interfaz WPF con un temporizador de 50 ms (`CemuFollower`) que intentaba forzar el anidamiento de ventanas Win32 sobre Cemu. En entornos Wayland, XWayland y gestores tipo mosaico (como Hyprland, Sway o SteamOS Game Mode), esto producía parpadeos continuos e inestabilidad de foco.

Esta suite resuelve el problema mediante la integración de **`MilkBar.CLI`**, un cliente mod headless desarrollado en [.NET 8](https://github.com/CarlosEvCode/MilkBarLauncher):
* **Cero Sobrecarga Gráfica:** La inyección de memoria, la comunicación por Named Pipes y la gestión de procesos de Cemu se ejecutan en segundo plano.
* **Nativo para Wayland y SteamOS:** Sin conflictos de anidamiento de ventanas Win32, sin robo de foco en XWayland y con estabilidad absoluta.

Regla opcional en Hyprland para Cemu:
```lua
o.window({ class = "^(cemu\\.exe)$" }, {
  opaque = true,
  suppress_event = "maximize",
})
```

---

## 5. Arquitectura Técnica y Notas de Compatibilidad

Como referencia técnica, se implementaron las siguientes soluciones para garantizar la estabilidad bajo Wine en Linux:

1. **Cliente Mod Headless (`MilkBar.CLI.exe`):**
   * Inicia Cemu, procesa los parámetros de configuración, inicializa los Named Pipes e inyecta `InjectDLL.dll` sin dependencias gráficas WPF.
2. **Wineprefix Compartido de 64 bits (`~/.local/share/wineprefixes/botw-multiplayer`):**
   * Cemu, el Servidor Dedicado y `MilkBar.CLI` se ejecutan en el mismo prefijo para que `InjectDLL.dll` y los Named Pipes puedan comunicarse con la memoria de Cemu sin bloqueos de pantalla blanca.
3. **Microsoft .NET 8.0 Desktop Runtime x64:**
   * Instalado directamente en el Wineprefix para solucionar cierres con código de salida 150.
4. **Validación de Paquetes de Modelos BCML:**
   * Incluye los 32 paquetes de modelos de Link (`Jugador1.sbactorpack` a `Jugador32.sbactorpack`) en `store_dir/merged/content/Actor/Pack/` enlazados a `BreathOfTheWild_BCML`.
5. **Enlace de Socket del Servidor:**
   * `ServerConfig.ini` se configura con `IP=127.0.0.1` y `DefaultGamemode=True` para enlazar inmediatamente el puerto `5050` y evitar excepciones de resolución DNS de `localhost`.
6. **Sanitización de Recursos en AppData Roaming:**
   * Extrae `QuestFlagsNames.txt` y `ArmorMapping.txt` en `%APPDATA%/BOTWM/` para evitar errores `FileNotFoundException`.

---

## 6. Créditos y Atribuciones

* **Milk Bar Launcher & Dedicated Server:** Desarrollado por la comunidad de [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher).
* **MilkBar.CLI (Cliente Headless):** Mantenido y adaptado para Linux por [CarlosEvCode](https://github.com/CarlosEvCode/MilkBarLauncher).
* **Cemu (Emulador de Wii U):** Desarrollado por [Team Cemu](https://cemu.info/) bajo licencia Mozilla Public License 2.0.
* **BCML:** Creado por [NiceneNerd](https://github.com/NiceneNerd/BCML).
* **Linux Manager y Despliegue:** Mantenido por [CarlosEvCode](https://github.com/CarlosEvCode).

### Descargo de Responsabilidad
Este repositorio no contiene ni distribuye archivos protegidos por derechos de autor de Nintendo, ejecutables oficiales (`.rpx`), claves de cifrado ni contenido propietario. Los usuarios deben suministrar su propia copia legal de *The Legend of Zelda: Breath of the Wild*, sus actualizaciones y DLC.
