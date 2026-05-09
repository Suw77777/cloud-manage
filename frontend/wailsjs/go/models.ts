export namespace main {
	
	export class InstanceView {
	    instanceId: string;
	    instanceName: string;
	    status: string;
	    regionId: string;
	    zoneId: string;
	    publicIp: string;
	    privateIp: string;
	    creationTime: string;
	
	    static createFrom(source: any = {}) {
	        return new InstanceView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.instanceId = source["instanceId"];
	        this.instanceName = source["instanceName"];
	        this.status = source["status"];
	        this.regionId = source["regionId"];
	        this.zoneId = source["zoneId"];
	        this.publicIp = source["publicIp"];
	        this.privateIp = source["privateIp"];
	        this.creationTime = source["creationTime"];
	    }
	}
	export class RegionResultView {
	    region: string;
	    instances: InstanceView[];
	    totalCount: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new RegionResultView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.region = source["region"];
	        this.instances = this.convertValues(source["instances"], InstanceView);
	        this.totalCount = source["totalCount"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MultiRegionResult {
	    success: boolean;
	    message: string;
	    regions: RegionResultView[];
	
	    static createFrom(source: any = {}) {
	        return new MultiRegionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.regions = this.convertValues(source["regions"], RegionResultView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OperationResult {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class QueryECSResult {
	    success: boolean;
	    message: string;
	    instances: InstanceView[];
	    totalCount: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryECSResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.instances = this.convertValues(source["instances"], InstanceView);
	        this.totalCount = source["totalCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

