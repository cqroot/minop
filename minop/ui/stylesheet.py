def get_stylesheet() -> str:
    colors: dict[str, str] = {
        "sidebar-background-color": "#2e3235",
        "sidebar-hover-color": "#464646",
        "foreground-color": "#1d1e21",
        "foreground-color-2": "#dddddd",
        "background-color": "#ffffff",
        "background-color-2": "#f6f6f6",
        "border-color": "#d9d9d9",
        "border-color2": "#5b5b5d",
        "hover-color": "#cacbcd",
        "checked-color": "#babbbc",
    }

    return """
/* ===== MWidget ===== */
MWidget {{
    background-color: {background-color};
}}

/* ===== MLabel ===== */
MLabel {{
    background-color: {background-color};
    color: {foreground-color};
}}

/* ===== MButton ===== */
MButton {{
    color: {foreground-color};
    background-color: {background-color};
    padding: 5px;
    border: none;
    border-radius: 5px;
}}

MButton:hover {{
    background-color: {hover-color};
}}

MButton:pressed {{
    background-color: {checked-color};
}}

/* ===== MInput ===== */
MInput {{
    background-color: {background-color};
    border: 1px solid {border-color};
    border-radius: 5px;
    padding: 3px;
}}

MInput:hover {{
    border-bottom: 1px solid {border-color2};
}}

MInput:focus {{
    background-color: {background-color-2};
    border-bottom: 1px solid {border-color2};
}}


/* ===== MSidebar ===== */
MSidebar {{
    outline: none;
    background-color: {sidebar-background-color};
    border: transparent;
    padding: 8px;
}}

MSidebar::item {{
    color: {foreground-color-2};
    background-color: {sidebar-background-color};
    border: transparent;
    border-radius: 5px;
    padding: 8px;
}}

MSidebar::item:selected {{
    color: {background-color};
    background-color: {sidebar-hover-color};
}}

/* ===== MOutput ===== */
MinopOutput QTextEdit {{
    border: 1px solid {border-color};
    border-radius: 5px;
    padding: 10px;
    font-family: "Cascadia Code", Consolas;
}}
""".format(
        **colors
    )
