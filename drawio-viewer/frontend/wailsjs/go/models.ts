export namespace models {
	
	export class Diagram {
	    id: string;
	    name: string;
	    xmlPath: string;
	    isTemp: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Diagram(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.xmlPath = source["xmlPath"];
	        this.isTemp = source["isTemp"];
	    }
	}

}

