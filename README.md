# resoluciones_mid_v2
:heavy_check_mark: Check: Repositorio API MID para el sistema de resoluciones v2.


## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)


### Variables de Entorno
```shell
# parametros de api
RESOLUCIONES_MID_V2_HTTP_PORT=[Puerto de exposición del API]
...
```
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf y están identificadas acorde a los lineamientos de definición de variables.V2...


### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/resoluciones_mid_v2

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/resoluciones_mid_v2

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
RESOLUCIONES_MID_V2_HTTP_PORT=8080 RESOLUCIONES_MID_V2_SOME_VARIABLE=valor bee run
```

### Ejecución Dockerfile
```shell
# Implementado para despliegue del Sistema de integración continua CI.
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/resoluciones_mid_v2

#2. Moverse a la carpeta del repositorio
cd resoluciones_mid_v2

#3. Crear un fichero con el nombre **custom.env**
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Ejecución Pruebas
```shell
# En Proceso
```

Pruebas unitarias
```shell
# En Proceso
```
## Estado CI


| Develop | Relese 0.0.1 | Master | Sonar |
| -- | -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/resoluciones_mid_v2/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/resoluciones_mid_v2) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/resoluciones_mid_v2/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/resoluciones_mid_v2) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/resoluciones_mid_v2/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/resoluciones_mid_v2) | [![Quality Gate Status](https://sonarqube.portaloas.udistrital.edu.co/api/project_badges/measure?project=udistrital%3Aresoluciones_mid_v2&metric=alert_status)](https://sonarqube.portaloas.udistrital.edu.co/dashboard?id=udistrital%3Aresoluciones_mid_v2) |


## Licencia

This file is part of resoluciones_docentes_mid.

resoluciones_docentes_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

resoluciones_docentes_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with resoluciones_docentes_mid. If not, see https://www.gnu.org/licenses/.
