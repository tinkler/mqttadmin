// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { User } from './user.model';


import { _HttpClient } from '@delon/theme';
import { modelUrlPrefix  } from './const';

export interface Page {
	
	page: number;
	
	perPage: number;
	
	total: number;
	
	
	fetchUser(): Promise<User[]> ;
	
}

export class Page {
	
	page: number = 0;
	
	perPage: number = 0;
	
	total: number = 0;
	

	constructor(
		private http: _HttpClient,
	){}

	
	fetchUser(): Promise<User[]>  {
		return new Promise((resolve, reject) => {
			this.http.post(`${modelUrlPrefix}/page/page/fetch-user`, { data: this, args: {  } }).subscribe({
				next: (res: { code: number; data: { data: any, resp: any }, message: string } ) => {
					if (res.code === 0) {
						this.page = res.data.data['page'];
						this.perPage = res.data.data['per_page'];
						this.total = res.data.data['total'];
						
						
						const resp = res.data.resp
						resolve(resp);
						
					} else {
						reject(res.message);
					}
				}, error: (err) => {
					reject(err);
				}
			});
		});
	}
	
}




