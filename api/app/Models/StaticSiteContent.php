<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class staticSiteContent extends Model
{
    use HasFactory;

    protected $hidden = [
        'id',
        'created_at',
        'updated_at'
    ];

    protected $fillable = [
        'lang_code',
        'content'
    ];
}
