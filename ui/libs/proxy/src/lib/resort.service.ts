import {Injectable} from '@angular/core';
import {from, mergeMap, Observable} from "rxjs";
import {ResortsEntity} from "@ui/data";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class ResortService {

  constructor(private httpClient: HttpClient) {
  }

  $loadResorts(): Observable<ResortsEntity> {
    return this.httpClient
      .get<ResortsEntity[]>("http://localhost:8080/resort")
      .pipe(mergeMap((value, index) => {
        return from(value)
      }))
  }
}
