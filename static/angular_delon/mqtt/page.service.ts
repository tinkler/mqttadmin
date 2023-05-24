// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
import { User } from './user.model';

import { Injectable } from '@angular/core';
import { Page,  } from './page.model';
import { UserService  } from './user.service';

import { _HttpClient } from '@delon/theme';

@Injectable({
	providedIn: 'root'
})
export class PageService {
  
	constructor(
		private http: _HttpClient,
		private userSrv: UserService,
		
		) { }

	
	newPage(): Page {
		return new Page(this.http);
	}
}
