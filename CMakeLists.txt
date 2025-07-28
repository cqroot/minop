cmake_minimum_required(VERSION 3.16)
project(minop LANGUAGES CXX)

find_package(Qt6 REQUIRED COMPONENTS Core Widgets Gui)

set(CMAKE_AUTOMOC ON)

include_directories(src)

add_executable(minop
    src/main.cpp
    src/main_window.cpp
    src/hosts_dock.cpp
    src/tasks_dock.cpp
)

target_link_libraries(minop
    Qt6::Core
    Qt6::Widgets
    Qt6::Gui
)

