export namespace main {
	
	export class CMSMetricsView {
	    instanceId: string;
	    cpuUtilization?: number;
	    memoryUtilization?: number;
	    diskReadBps?: number;
	    diskWriteBps?: number;
	    internetRx?: number;
	    internetTx?: number;
	    updateTime: string;
	
	    static createFrom(source: any = {}) {
	        return new CMSMetricsView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.instanceId = source["instanceId"];
	        this.cpuUtilization = source["cpuUtilization"];
	        this.memoryUtilization = source["memoryUtilization"];
	        this.diskReadBps = source["diskReadBps"];
	        this.diskWriteBps = source["diskWriteBps"];
	        this.internetRx = source["internetRx"];
	        this.internetTx = source["internetTx"];
	        this.updateTime = source["updateTime"];
	    }
	}
	export class CMSMetricsResult {
	    success: boolean;
	    message: string;
	    metrics: CMSMetricsView[];
	
	    static createFrom(source: any = {}) {
	        return new CMSMetricsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.metrics = this.convertValues(source["metrics"], CMSMetricsView);
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
	export class LogEntryView {
	    timestamp: number;
	    content: Record<string, string>;

	    static createFrom(source: any = {}) {
	        return new LogEntryView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.content = source["content"];
	    }
	}
	export class SLSLogStoreResult {
	    success: boolean;
	    message: string;
	    logStores: string[];

	    static createFrom(source: any = {}) {
	        return new SLSLogStoreResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.logStores = source["logStores"];
	    }
	}
	export class SLSLogQueryResult {
	    success: boolean;
	    message: string;
	    logs: LogEntryView[];
	    count: number;
	    hasMore: boolean;

	    static createFrom(source: any = {}) {
	        return new SLSLogQueryResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.logs = this.convertValues(source["logs"], LogEntryView);
	        this.count = source["count"];
	        this.hasMore = source["hasMore"];
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
	export class SLSStreamResult {
	    success: boolean;
	    message: string;

	    static createFrom(source: any = {}) {
	        return new SLSStreamResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class BucketView {
	    name: string;
	    location: string;
	    creationDate: string;
	    storageClass: string;
	    extranetEndpoint: string;
	    intranetEndpoint: string;

	    static createFrom(source: any = {}) {
	        return new BucketView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.location = source["location"];
	        this.creationDate = source["creationDate"];
	        this.storageClass = source["storageClass"];
	        this.extranetEndpoint = source["extranetEndpoint"];
	        this.intranetEndpoint = source["intranetEndpoint"];
	    }
	}
	export class ObjectView {
	    key: string;
	    size: number;
	    lastModified: string;
	    etag: string;
	    type: string;
	    storageClass: string;
	    isFolder: boolean;

	    static createFrom(source: any = {}) {
	        return new ObjectView(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.size = source["size"];
	        this.lastModified = source["lastModified"];
	        this.etag = source["etag"];
	        this.type = source["type"];
	        this.storageClass = source["storageClass"];
	        this.isFolder = source["isFolder"];
	    }
	}
	export class OSSBucketResult {
	    success: boolean;
	    message: string;
	    buckets: BucketView[];

	    static createFrom(source: any = {}) {
	        return new OSSBucketResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.buckets = this.convertValues(source["buckets"], BucketView);
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
	export class OSSObjectResult {
	    success: boolean;
	    message: string;
	    objects: ObjectView[];
	    isTruncated: boolean;

	    static createFrom(source: any = {}) {
	        return new OSSObjectResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.objects = this.convertValues(source["objects"], ObjectView);
	        this.isTruncated = source["isTruncated"];
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

