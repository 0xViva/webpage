class Wiigen extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        const shadow = this.attachShadow({ mode: "open" });
        const anchor = document.createElement("a");
        anchor.setAttribute("class", "bagde")
        anchor.setAttribute("href", "https://new.jpwiig.dev")
        const l = document.createElement("div");
        l.setAttribute("class", "l")
        l.textContent = "λ~/";
        anchor.setAttribute("target", "_blank");
        anchor.setAttribute("aria-label", "Bagde for jpwiig.dev")
        anchor.appendChild(l);
        const sep = document.createElement("div")
        sep.setAttribute("class", "sep")
        anchor.appendChild(sep)
        const r = document.createElement("div");
        r.textContent = "jpwiig.dev";
        r.setAttribute("class", "r")
        anchor.appendChild(r);
        const style = document.createElement("style");
        style.textContent = `
        .bagde {
            display: flex;
            width: 210px;
            height: 30px;
            border: 3px groove green;
            box-sizing: border-box;
            font-size: 1rem;
            text-decoration: none;
            outline: none; 
  
            .l, .r {
                display: flex;
                justify-content: center;
                align-items: center;
                margin: 1px;
                height: calc(100% - 2px);
                font-family: 'monospace'; 
            }
            

                        .l{
            background:  rgb(4, 228, 4);
            color: black; 
            padding-right: 5px;  
            
            }
            .r{
            background: black; 
            text-transform:uppercase;
            color:rgb(4, 228, 4) ;
            outline: none;
            }
            .sep {
            width: 5px;
            background-color: black;
            }
        }`;


        shadow.appendChild(style);
        shadow.appendChild(anchor);
    }

}
customElements.define("wiigen-bagde", Wiigen);
