class FrkFrontendBadge extends HTMLElement {
    constructor() {
        super();
    }

    async connectedCallback() {
        const shadow = this.attachShadow({
            mode: "open",
        });

        const anchor = document.createElement("a");
        anchor.setAttribute("class", "badge");
        anchor.setAttribute("href", "https://frkfrontend.no");
        anchor.setAttribute("target", "_blank");
        anchor.setAttribute(
            "aria-label",
            "A badge to get you to Elises website frkfrontend.no"
        );

        const left = document.createElement("div");
        left.setAttribute("class", "left");
        left.textContent = "🦄";
        anchor.appendChild(left);

        const separator = document.createElement("div");
        separator.setAttribute("class", "separator");
        anchor.appendChild(separator);

        const right = document.createElement("div");
        right.setAttribute("class", "right");
        right.textContent = "FRK. FRONTEND";
        anchor.appendChild(right);

        const style = document.createElement("style");
        style.textContent = `
        .badge {
            display: flex;
            width: 210px;
            height: 30px;
            border: 1px solid black;
            box-sizing: border-box;
            font-size: 0.7rem;
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
                font-family: 'Press Start 2P';
            }
    
            .left {
                flex: 1;
                background-color: rgb(255, 224, 230);
                padding: 0 1px;
            }

            .separator {
                width: 1px;
                background-color: black;
            }

            .right {
                flex: 5;
                background-color:rgb(158 84 241);
                color: rgb(255, 224, 230);
                text-transform: uppercase;
            }
        }
        `;

        shadow.appendChild(style);
        shadow.appendChild(anchor);
    }
}

customElements.define("frk-frontend-badge", FrkFrontendBadge);
