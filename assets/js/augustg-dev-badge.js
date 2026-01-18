class AugustGDevBadge extends HTMLElement {
    constructor() {
        super();
    }

    async connectedCallback() {
        const shadow = this.attachShadow({
            mode: "open"
        });

        const anchor = document.createElement("a");
        anchor.setAttribute("class", "badge");
        anchor.setAttribute("href", "https://augustg.dev");
        anchor.setAttribute("target", "_blank");
        anchor.setAttribute("aria-label", "A badge representing AugustG's Webpage");

        const left = document.createElement("div");
        left.setAttribute("class", "left");
        left.textContent = "🐸";
        anchor.appendChild(left);

        const separator = document.createElement("div");
        separator.setAttribute("class", "separator");
        anchor.appendChild(separator);

        const right = document.createElement("div");
        right.setAttribute("class", "right");
        right.textContent = "augustg.dev";
        anchor.appendChild(right);

        const style = document.createElement("style");
        style.textContent = `
        .badge {
            display: flex;
            width: 210px;
            height: 30px;
            border: 1px solid blue;
            box-sizing: border-box;
            font-size: 1.0rem;
            text-decoration: none;

            &:active {
                color: initial;
            }

            .left, .right {
                display: flex;
                justify-content: center;
                align-items: center;
                margin: 1px;
                height: calc(100% - 2px);
                font-family: 'JetBrainsMono Nerd';
            }
    
            .left {
                flex: 1;
                background-color: #0065ff;
                padding: 0 1px;
            }

            .separator {
                width: 1px;
                background-color: blue;
            }

            .right {
                flex: 5;
                background-color: #2f3138;

                color: #01a8f4;
                text-transform: uppercase;
            }
        }
        `;

        shadow.appendChild(style);
        shadow.appendChild(anchor);
    }
}

customElements.define('augustg-dev-badge', AugustGDevBadge);
